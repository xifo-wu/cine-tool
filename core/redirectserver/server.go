package redirectserver

import (
	"cine-tool/core/redirectserver/alist"
	"cine-tool/core/redirectserver/emby"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func Run() {
	PORT := viper.GetString("PORT_302")

	e := echo.New()
	e.HideBanner = true

	embyURL, _ := url.Parse(viper.GetString("EMBY_URL"))
	proxy := httputil.NewSingleHostReverseProxy(embyURL)

	e.Any("/*actions", func(c echo.Context) error {
		currentURI := c.Request().RequestURI
		videoID, err := extractIDFromPath(currentURI)
		if err != nil {
			proxy.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
		log.Println("【EMBY 302 服务】Request URI:", currentURI)
		log.Println("【EMBY 302 服务】Header:", c.Request().Header)

		mediaSourceID := c.QueryParam("MediaSourceId")
		if mediaSourceID == "" {
			mediaSourceID = c.QueryParam("mediaSourceId")
		}

		if videoID == "" || mediaSourceID == "" {
			proxy.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}

		itemInfoUri, itemId, etag, mediaSourceId, apiKey := emby.GetItemPathInfo(c)
		embyRes, err := emby.GetEmbyItems(itemInfoUri, itemId, etag, mediaSourceId, apiKey)

		if err != nil {
			e.Logger.Error(fmt.Sprintf("获取 Emby 失败。错误信息: %v", err))
			proxy.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}

		log.Println("Emby 原地址：" + embyRes["path"].(string))
		e.Logger.Info("Emby 原地址：" + embyRes["path"].(string))
		alistPath := replacePath(embyRes["path"].(string))
		alistPath = ensureLeadingSlash(alistPath)

		originalHeaders := make(map[string]string)
		for key, value := range c.Request().Header {
			if len(value) > 0 {
				originalHeaders[key] = value[0]
			}
		}

		sign := alist.Sign(alistPath, 0)
		alistFullUrl := viper.GetString("ALIST_URL") + "/d" + alistPath + "?sign=" + sign
		e.Logger.Info("Alist 原地址：" + alistFullUrl)

		redirectURL, err := alist.GetRedirectURL(alistFullUrl, originalHeaders)
		if err != nil {
			c.Logger().Error(fmt.Sprintf("获取 Alist 地址失败。错误信息: %v", err))
			proxy.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}

		e.Logger.Info("获取重定向链接：" + redirectURL)
		log.Println("获取重定向链接：" + redirectURL)

		return c.Redirect(302, redirectURL)
	})

	go func() {
		e.Start(":" + PORT)
	}()
}

func extractIDFromPath(path string) (string, error) {
	// 编译正则表达式
	re := regexp.MustCompile(`/[Vv]ideos/(\S+)/(stream|original|master)`)
	// 执行匹配操作
	matches := re.FindStringSubmatch(path)

	// 如果找到匹配项，第一个分组就是我们想要的视频ID
	if len(matches) >= 2 {
		return matches[1], nil
	}

	// 如果没有匹配项，返回错误
	return "", fmt.Errorf("no match found")
}

func ensureLeadingSlash(alistPath string) string {
	if !strings.HasPrefix(alistPath, "/") {
		alistPath = "/" + alistPath // 不是以 / 开头，加上 /
	}

	alistPath = convertToLinuxPath(alistPath)
	return alistPath
}

func convertToLinuxPath(windowsPath string) string {
	// 将所有的反斜杠转换成正斜杠
	linuxPath := strings.ReplaceAll(windowsPath, "\\", "/")
	return linuxPath
}

func replacePath(path string) string {
	replaces := viper.GetStringSlice("REPLACE_PATHS")
	for key, value := range replaces {
		log.Println(strings.HasPrefix(path, value), path, key, value, "key, value 没匹配到？？？")
		if strings.HasPrefix(path, value) {
			return strings.Replace(path, value, "", 1)
		}
	}
	return path
}
