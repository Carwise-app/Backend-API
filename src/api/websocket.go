package main

import (
	"carwise"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id         string
	UserId     string
	ListingId  string
	ReceiverId string
	Conn       *websocket.Conn
	Send       chan []byte
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

type Message struct {
	Message    string `json:"message"`
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	ListingId  string `json:"listing_id"`
	Timestamp  int64  `json:"timestamp"`
}

var hub = &Hub{
	Clients:    make(map[*Client]bool),
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s", client.Id)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: %s", client.Id)

		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					delete(h.Clients, client)
					close(client.Send)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		isReceiverActive := false
		hub.mu.RLock()
		for client := range hub.Clients {
			if client.UserId == msg.ReceiverId && client.ListingId == msg.ListingId {
				isReceiverActive = true
				break
			}
		}
		hub.mu.RUnlock()

		interactor.SendMessage(&carwise.SendMessageRequest{
			Message:          msg.Message,
			UserId:           msg.SenderId,
			ReceiverId:       msg.ReceiverId,
			ListingId:        msg.ListingId,
			IsReceiverActive: isReceiverActive,
		})

		c.broadcastToRoom(&msg)
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}

	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c *Client) broadcastToRoom(msg *Message) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	for client := range hub.Clients {
		if client.ListingId == c.ListingId &&
			client.UserId == msg.ReceiverId {
			msg.Timestamp = time.Now().Unix()
			msgBytes, _ := json.Marshal(msg)
			select {
			case client.Send <- msgBytes:
			default:
				delete(hub.Clients, client)
				close(client.Send)
			}
		}
	}
}

func generateClientID(userID, listingID, receiverID string) string {
	return fmt.Sprintf("%s_%s_%s", userID, listingID, receiverID)
}
