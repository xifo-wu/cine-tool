package app

import (
	"cine-tool/app/model"
	"cine-tool/app/utils/watcher"
	"cine-tool/core"
	"cine-tool/core/redirectserver"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

func RunServer() {
	PORT := viper.GetString("PORT")
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	db := core.InitDB()
	db.AutoMigrate(model.CloudSymlinkSync{})
	w, err := InitWatcher()
	if err != nil {
		e.Logger.Fatal(err)
	}
	InitRouter(e, db, w)

	e.Static("/", "dist")

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: nil,
		// Root directory from where the static content is served.
		Root: "dist",
		// Index file for serving a directory.
		// Optional. Default value "index.html".
		Index: "index.html",
		// Enable HTML5 mode by forwarding all not-found requests to root so that
		// SPA (single-page application) can handle the routing.
		HTML5:      true,
		Browse:     false,
		IgnoreBase: false,
		Filesystem: nil,
	}))

	syncWatcher := &watcher.SyncWatcher{
		Watcher: w,
	}

	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				syncWatcher.HandleEvent(event)

			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Errorf("error:", err)
			}
		}
	}()

	defer w.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ENABLE_302 := viper.GetBool("ENABLE_302")
	if ENABLE_302 {
		redirectserver.Run()
	}

	go func() {
		if err := e.Start(":" + PORT); err != nil {
			e.Logger.Info("Shutting Down The Server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
