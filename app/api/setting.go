package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type SettingApi struct {
}

func (api *SettingApi) List(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"avatarUrl": viper.GetString("AVATAR_URL"),
		"nickname":  viper.GetString("NICKNAME"),
		"version":   "0.0.1.beta",
	})
}
