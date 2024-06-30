package api

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type LoginApi struct {
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *LoginApi) Login(c echo.Context) error {
	var data LoginBody
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "message": err.Error()})
	}

	// Throws unauthorized error
	if data.Username != viper.GetString("USER") || data.Password != viper.GetString("PASS") {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		viper.GetString("USER"),
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(viper.GetInt("LoginDuration")))),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(viper.GetString("USER")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
