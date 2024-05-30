package controllers

import (
	"github.com/HakimHC/altostratus-golang-auth/models"
	"github.com/HakimHC/altostratus-golang-auth/responses"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Register(c echo.Context) error {
	body := new(models.AuthRequestDTO)

	if err := c.Bind(body); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	userNameExists, err := getUserByUsername(body.UserName, "Users")
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else if userNameExists != nil {
		return ErrorResponse(c, http.StatusForbidden, "Username is already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password.")
	}

	u := models.User{
		ID:        uuid.New().String(),
		Username:  body.UserName,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := putItemInDynamoDB(u, "Users"); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	claims := &models.JwtCustomClaims{
		UserId: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// TODO: get secret from env
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated,
		responses.AuthResponse{
			Status:  http.StatusCreated,
			Message: "created",
			Data: &echo.Map{
				"user_name":    u.Username,
				"created_at":   u.CreatedAt,
				"access_token": t,
			},
		},
	)
}

func Login(c echo.Context) error {
	body := new(models.AuthRequestDTO)

	if err := c.Bind(body); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(body); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := getUserByUsername(body.UserName, "Users")
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve user")
	}
	if user == nil {
		return ErrorResponse(c, http.StatusUnauthorized, "Username or password incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return ErrorResponse(c, http.StatusUnauthorized, "Username or password incorrect")
	}

	claims := &models.JwtCustomClaims{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Token expires in 72 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// TODO: get it from env
	//secret := os.Getenv("JWT_SECRET")
	secret := "secret"
	if secret == "" {
		return ErrorResponse(c, http.StatusInternalServerError, "JWT secret is not set.")
	}

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "Failed to sign token")
	}

	return c.JSON(
		http.StatusCreated,
		responses.AuthResponse{
			Status:  http.StatusOK,
			Message: "ok",
			Data: &echo.Map{
				"access_token": t,
			},
		})
}
