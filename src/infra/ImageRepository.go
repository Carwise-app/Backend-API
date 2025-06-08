package infra

import (
	"carwise"
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type ImageRepository struct {
	db       *sql.DB
	basePath string
}

func NewImageRepository() *ImageRepository {
	return &ImageRepository{
		db:       ConnectDb(),
		basePath: "./uploads",
	}
}

func (r *ImageRepository) SaveImage(file *multipart.FileHeader, userId string) (*carwise.Image, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	imageId := uuid.New().String()

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", imageId, ext)

	dirPath := filepath.Join(r.basePath, userId)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(dirPath, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	image := &carwise.Image{
		Id:        imageId,
		Path:      fmt.Sprintf("%s/%s/%s", r.basePath, userId, filename),
		CreatedBy: userId,
		CreatedAt: time.Now().Unix(),
	}

	query := `
		INSERT INTO images (id, path, created_by, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err = r.db.QueryRow(query, image.Id, image.Path, image.CreatedBy, image.CreatedAt).Scan(&image.Id)
	if err != nil {

		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save image record to database: %w", err)
	}

	return image, nil
}

func (r *ImageRepository) GetImageById(id string) (*carwise.Image, error) {
	query := `
		SELECT id, path, created_by, created_at
		FROM images
		WHERE id = $1
	`

	var image carwise.Image
	err := r.db.QueryRow(query, id).Scan(&image.Id, &image.Path, &image.CreatedBy, &image.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return &image, nil
}

func (r *ImageRepository) DeleteImage(id string) error {
	log.Println("Deleting image with ID:", id)
	image, err := r.GetImageById(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM images WHERE id = $1`
	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete image from database: %w", err)
	}

	if err := os.Remove(image.Path); err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}

func (r *ImageRepository) ImagesCount() (int, error) {
	query := `
		SELECT COUNT(*) FROM images
	`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}
