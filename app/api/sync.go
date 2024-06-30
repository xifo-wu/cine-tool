package api

import (
	"cine-tool/app/model"
	"cine-tool/app/utils/symlink"
	"cine-tool/app/utils/watcher"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SyncApi struct {
	Api
	SyncWatcher *watcher.SyncWatcher
}

// Create 创建一个同步软链接目录
// 1. 添加目录监控
// 2. 递归创建软链接
func (s *SyncApi) Create(c echo.Context) error {
	var cloudSymlinkSync model.CloudSymlinkSync
	if err := c.Bind(&cloudSymlinkSync); err != nil {
		c.Logger().Error(err)
		return c.JSON(400, map[string]any{"success": false, "message": "Invalid request"})
	}
	// TODO 得判断是否是别的路径的子目录，是的话不允许创建

	s.DB.Create(&cloudSymlinkSync)
	if cloudSymlinkSync.ID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "error": "创建失败"})
	}

	go func() {
		// 添加监控
		s.SyncWatcher.AddSyncDirs(cloudSymlinkSync.CloudPath)
		// 创建软链接
		symlink.CreateSymlink(cloudSymlinkSync.CloudPath, cloudSymlinkSync.LocalPath)
	}()

	return c.JSON(http.StatusOK, map[string]any{"success": true, "data": cloudSymlinkSync})
}

func (s *SyncApi) Get(e echo.Context) error {
	return nil
}

func (s *SyncApi) List(c echo.Context) error {
	var cloudSymlinkSyncs []model.CloudSymlinkSync
	s.DB.Find(&cloudSymlinkSyncs)
	return c.JSON(http.StatusOK, map[string]any{"success": true, "data": cloudSymlinkSyncs})
}

func (s *SyncApi) Update(c echo.Context) error {
	return nil
}

func (s *SyncApi) Delete(c echo.Context) error {
	var cloudSymlinkSync model.CloudSymlinkSync
	id := c.Param("id")

	log.Println("delete id:", id)

	s.DB.First(&cloudSymlinkSync, id)
	if cloudSymlinkSync.ID == 0 {
		return c.JSON(http.StatusOK, map[string]any{"success": false, "message": "未找到该记录"})
	}

	go func() {
		// TODO 添加一个参数是否删除本地文件夹。
		s.SyncWatcher.RemoveSyncDirs(cloudSymlinkSync.CloudPath)
	}()

	err := s.DB.Delete(&cloudSymlinkSync).Error
	if err != nil {
		c.Logger().Errorf("删除监控目录失败: %s", cloudSymlinkSync.CloudPath)
		return c.JSON(http.StatusOK, map[string]any{"success": false, "message": "删除失败"})
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true, "message": "删除成功"})
}

func (s *SyncApi) ReGenerate(c echo.Context) error {
	return nil
}
