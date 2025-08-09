package controllers

import (
	"net/http"
	"strconv"

	"github.com/blog-platform/domain"
	"github.com/gin-gonic/gin"
)

type UserRegisterDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginDTO struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type ResetPasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email"`
}

type UpdatePasswordDirectDTO struct {
	NewPassword string `json:"new_password"`
}

type UserController struct {
	userUsecase domain.IUserUsecase
}

func NewUserController(uu domain.IUserUsecase) *UserController {
	return &UserController{
		userUsecase: uu,
	}
}

func (uc *UserController) Register(ctx *gin.Context) {
	var userInput UserRegisterDTO

	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user := domain.User{
		Email:    userInput.Email,
		Username: userInput.Username,
		Password: userInput.Password,
	}

	user, err := uc.userUsecase.Register(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func (uc *UserController) ActivateAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	err := uc.userUsecase.ActivateAccount(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user account activated"})
}

func (uc *UserController) Login(ctx *gin.Context) {
	var userInput UserLoginDTO

	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	accessToken, refreshToken, err := uc.userUsecase.Login(userInput.Identifier, userInput.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access":  accessToken,
		"refresh": refreshToken,
		"message": "Logged in successfully",
	})
}

func (uc *UserController) GetProfile(ctx *gin.Context) {
	idParam := ctx.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	user, err := uc.userUsecase.GetUserProfile(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) Promote(ctx *gin.Context) {
	id := ctx.Param("id")
	err := uc.userUsecase.Promote(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}

func (uc *UserController) Demote(ctx *gin.Context) {
	id := ctx.Param("id")
	err := uc.userUsecase.Demote(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user demoted to user"})
}
  
func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	idParam := ctx.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	err = uc.userUsecase.UpdateUserProfile(userID, updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

func (uc *UserController) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	access, refresh, err := uc.userUsecase.RefreshToken(authHeader)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"access": access, "refresh": refresh})
}

func (uc *UserController) ResetPassword(ctx *gin.Context) {
	var body ResetPasswordDTO
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, _ := userIDVal.(string)
	if err := uc.userUsecase.ResetPassword(userID, body.OldPassword, body.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

func (uc *UserController) ForgotPassword(ctx *gin.Context) {
	var body ForgotPasswordDTO
	if err := ctx.ShouldBindJSON(&body); err != nil || body.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := uc.userUsecase.ForgotPassword(body.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "reset link sent (check email)"})
}

func (uc *UserController) UpdatePasswordDirect(ctx *gin.Context) {
	userID := ctx.Param("id")
	var body UpdatePasswordDirectDTO
	if err := ctx.ShouldBindJSON(&body); err != nil || body.NewPassword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token required"})
		return
	}
	if err := uc.userUsecase.UpdatePasswordDirect(userID, body.NewPassword, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "password updated"})
}
