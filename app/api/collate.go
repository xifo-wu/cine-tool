package api

import (
	"cine-tool/app/utils/collate"

	"github.com/labstack/echo/v4"
)

type CollateApi struct {
}

func (c *CollateApi) Collate(e echo.Context) error {
	var collateRequest collate.CollateRequest
	if err := e.Bind(&collateRequest); err != nil {
		e.Logger().Error(err)
		return e.JSON(400, map[string]any{"success": false, "message": "Invalid request"})
	}

	mediaInfo, err := collate.GetMediaInfo(collateRequest)
	if err != nil {
		return e.JSON(400, map[string]any{"success": false, "message": err.Error()})
	}

	return e.JSON(200, mediaInfo)
}
