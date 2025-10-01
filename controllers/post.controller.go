package controllers

import (
	"net/http"
	"strconv"

	"app/models"
	"app/services"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService *services.PostService
}

func NewPostController(postService *services.PostService) *PostController {
	return &PostController{postService: postService}
}

// CreatePost handles post creation
func (pc *PostController) CreatePost(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreatePostRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	newPost, err := pc.postService.CreatePost(payload, currentUser.ID)
	if err != nil {
		switch err.Error() {
		case "post with that title already exists":
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
		default:
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPost})
}

// UpdatePost handles post updates
func (pc *PostController) UpdatePost(ctx *gin.Context) {
	postID := ctx.Param("postId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdatePost
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedPost, err := pc.postService.UpdatePost(postID, payload, currentUser.ID)
	if err != nil {
		switch err.Error() {
		case "post not found":
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
		default:
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedPost})
}

// FindPostById handles retrieving a single post
func (pc *PostController) FindPostById(ctx *gin.Context) {
	postID := ctx.Param("postId")

	post, err := pc.postService.FindPostByID(postID)
	if err != nil {
		switch err.Error() {
		case "post not found":
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
		default:
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

// FindPosts handles retrieving paginated posts
func (pc *PostController) FindPosts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)
	if err != nil || intPage < 1 {
		intPage = 1
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil || intLimit < 1 {
		intLimit = 10
	}

	posts, err := pc.postService.FindPosts(intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "data": posts})
}

// DeletePost handles post deletion
func (pc *PostController) DeletePost(ctx *gin.Context) {
	postID := ctx.Param("postId")

	err := pc.postService.DeletePost(postID)
	if err != nil {
		switch err.Error() {
		case "post not found":
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
		default:
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
