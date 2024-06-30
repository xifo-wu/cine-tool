package collate

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"

	"github.com/nssteinbrenner/anitogo"
)

type FileInfo struct {
	EpisodeNumber    string
	Year             string
	ReleaseGroup     string
	AudioTerm        string
	VideoResolution  string
	Subtitles        string
	Language         string
	VideoTerm        string
	EpisodeNumberAlt string
	FileExtension    string
	Season           int
	FileTitle        string
	Title            string
}

func zeroPad(num int) string {
	if num < 10 {
		return fmt.Sprintf("0%d", num)
	}
	return fmt.Sprintf("%d", num)
}

func ParseFileName(mediaInfo MediaInfo, fileName string) string {
	parsed := anitogo.Parse(fileName, anitogo.DefaultOptions)

	log.Println(parsed, "parsed")

	fileInfo := FileInfo{}
	fileInfo.FileTitle = parsed.AnimeTitle
	fileInfo.Title = mediaInfo.Title
	fileInfo.EpisodeNumber = strings.Join(parsed.EpisodeNumber, "-")
	fileInfo.Year = parsed.AnimeYear
	fileInfo.ReleaseGroup = parsed.ReleaseGroup
	fileInfo.AudioTerm = strings.Join(parsed.AudioTerm, ".")
	fileInfo.VideoResolution = parsed.VideoResolution
	fileInfo.Subtitles = strings.Join(parsed.Subtitles, ".")
	fileInfo.Language = strings.Join(parsed.Language, ".")
	fileInfo.VideoTerm = strings.Join(parsed.VideoTerm, ".")
	fileInfo.EpisodeNumberAlt = strings.Join(parsed.EpisodeNumberAlt, ".")
	fileInfo.FileExtension = parsed.FileExtension

	var season int
	if len(parsed.AnimeSeason) > 0 {
		season, _ = strconv.Atoi(parsed.AnimeSeason[0])
	}

	fileInfo.Season = season
	if fileInfo.Season == 0 {
		fileInfo.Season, _ = strconv.Atoi(mediaInfo.Season)
	}

	tpl := "{{.Title}} - S{{zeroPad .Season}}E{{.EpisodeNumber}} - {{if .VideoResolution}}{{toUpper .VideoResolution}}.{{end}}{{if .Subtitles}}{{.Subtitles}}.{{end}}{{if .Language}}{{.Language}}.{{end}}{{if .VideoTerm}}{{.VideoTerm}}.{{end}}{{if .AudioTerm}}{{.AudioTerm}}.{{end}}{{.FileExtension}}"

	tmpl, err := template.New("filename").Funcs(template.FuncMap{
		"toUpper": strings.ToUpper,
		"zeroPad": zeroPad,
	}).Parse(tpl)
	if err != nil {
		return fileName
	}
	// 创建一个字节缓冲区用于存储格式化后的字符串
	var tplBuffer bytes.Buffer

	// 使用模板将数据填充到模板中，并存储到缓冲区中
	err = tmpl.Execute(&tplBuffer, fileInfo)
	if err != nil {
		panic(err)
	}

	// 从缓冲区中获取格式化后的字符串
	formattedString := tplBuffer.String()

	return formattedString
}
