package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostService struct {
	DB *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db}
}

// CreatePost creates a new post
func (s *PostService) CreatePost(payload *models.CreatePostRequest, userID uuid.UUID) (*models.Post, error) {
	now := time.Now()
	newPost := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		User:      userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := s.DB.Create(&newPost)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			return nil, errors.New("post with that title already exists")
		}
		return nil, fmt.Errorf("failed to create post: %w", result.Error)
	}

	return &newPost, nil
}

// UpdatePost updates an existing post
func (s *PostService) UpdatePost(postID string, payload *models.UpdatePost, userID uuid.UUID) (*models.Post, error) {
	var existingPost models.Post
	result := s.DB.First(&existingPost, "id = ?", postID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, fmt.Errorf("failed to fetch post: %w", result.Error)
	}

	// Optional: Add authorization check
	// if existingPost.User != userID {
	//     return nil, errors.New("unauthorized to update this post")
	// }

	now := time.Now()
	updatedData := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		User:      userID,
		UpdatedAt: now,
	}

	result = s.DB.Model(&existingPost).Updates(updatedData)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update post: %w", result.Error)
	}

	return &existingPost, nil
}

// FindPostByID retrieves a post by ID
func (s *PostService) FindPostByID(postID string) (*models.Post, error) {
	var post models.Post
	result := s.DB.First(&post, "id = ?", postID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, fmt.Errorf("failed to fetch post: %w", result.Error)
	}

	return &post, nil
}

// FindPosts retrieves a paginated list of posts
func (s *PostService) FindPosts(page, limit int) ([]models.Post, error) {
	offset := (page - 1) * limit

	var posts []models.Post
	result := s.DB.Limit(limit).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", result.Error)
	}

	return posts, nil
}

// DeletePost deletes a post by ID
func (s *PostService) DeletePost(postID string) error {
	result := s.DB.Delete(&models.Post{}, "id = ?", postID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete post: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("post not found")
	}

	return nil
}
