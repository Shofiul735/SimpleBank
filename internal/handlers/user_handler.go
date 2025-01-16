package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shofiul735/simple_bank/internal/core/domain"
	"github.com/shofiul735/simple_bank/internal/core/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// RegisterRoutes sets up the user routes
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/", h.createUser)    // POST /api/users
	r.GET("/:userID", h.getUser) // GET /api/users/:userID
}

// createUser handles user creation
func (h *UserHandler) createUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		switch err {
		case services.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case services.ErrDuplicateEmail:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// getUser handles fetching a single user
func (h *UserHandler) getUser(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}
