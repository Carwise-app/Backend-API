package infra

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type CDNRepository struct {
	basePath string
}

func NewCDNRepository() *CDNRepository {
	return &CDNRepository{
		basePath: "/images",
	}
}

func (r *CDNRepository) SaveUserAvatar(userID string, image io.Reader) (string, error) {
	dirPath := filepath.Join(r.basePath, "users", userID)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(dirPath, "avatar.png")

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, image)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/users/%s/avatar.png", r.basePath, userID), nil
}
