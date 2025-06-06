package carwise

type Notification struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	Status    int               `json:"status"`
	Data      map[string]string `json:"data"`
	Read      bool              `json:"read"`
	CreatedBy string            `json:"created_by"`
	CreatedAt int64             `json:"created_at"`
}
