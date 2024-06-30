package app

import (
	"cine-tool/app/api"
	"cine-tool/app/utils/watcher"
	"net/http"

	"github.com/fsnotify/fsnotify"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func InitRouter(e *echo.Echo, db *gorm.DB, w *fsnotify.Watcher) {
	collateApi := api.CollateApi{}

	apiGroup := e.Group("/api")
	settingApi := api.SettingApi{}
	settingGroupApi := apiGroup.Group("/settings")
	settingGroupApi.GET("", settingApi.List)

	loginApi := api.LoginApi{}
	apiGroup.POST("/login", loginApi.Login)

	auth := apiGroup.Group("")
	auth.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(viper.GetString("USER")),
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]any{"code": 401, "msg": "请先登录"})
		},
	}))

	auth.POST("/collate", collateApi.Collate)

	syncApi := api.SyncApi{
		Api: api.Api{DB: db},
		SyncWatcher: &watcher.SyncWatcher{
			Watcher: w,
		},
	}

	syncGroupApi := auth.Group("/sync")
	syncGroupApi.GET("", syncApi.List)
	syncGroupApi.POST("", syncApi.Create)
	syncGroupApi.PUT("/:id", syncApi.Update)
	syncGroupApi.GET("/:id", syncApi.Get)
	syncGroupApi.DELETE("/:id", syncApi.Delete)
}
