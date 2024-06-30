package collate

import (
	"fmt"
	"os"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/spf13/viper"
)

type CollateRequest struct {
	TmdbID      int    `json:"tmdbId"`
	MediaType   string `json:"mediaType"`
	CollatePath string `json:"collatePath"`
	TargetPath  string `json:"targetPath"`
	Season      string `json:"season"`
}

type MediaInfo struct {
	MediaType string            `json:"mediaType"`
	Title     string            `json:"title"`
	Season    string            `json:"season"`
	Year      string            `json:"year"`
	List      map[string]string `json:"list"`
}

func GetMediaInfo(collateRequest CollateRequest) (MediaInfo, error) {
	mediaInfo := MediaInfo{MediaType: collateRequest.MediaType, List: make(map[string]string), Season: collateRequest.Season}
	tmdbClient, err := tmdb.Init(viper.GetString("TMDB_API_KEY"))
	if err != nil {
		return mediaInfo, fmt.Errorf("error initializing tmdb client: %v", err)
	}

	options := make(map[string]string)
	options["language"] = "zh-CN"

	if collateRequest.MediaType == "movie" {
		tv, err := tmdbClient.GetMovieDetails(collateRequest.TmdbID, options)
		if err != nil {
			return mediaInfo, fmt.Errorf("error getting movie details")
		}
		mediaInfo.Title = tv.Title
		mediaInfo.Year = tv.ReleaseDate[0:4]
	}

	if collateRequest.MediaType == "tv" {
		movie, err := tmdbClient.GetTVDetails(collateRequest.TmdbID, options)
		if err != nil {
			return mediaInfo, fmt.Errorf("error getting tv details")
		}

		mediaInfo.Title = movie.Name
		mediaInfo.Year = movie.FirstAirDate[0:4]
	}

	files, err := os.ReadDir(collateRequest.CollatePath)
	if err != nil {
		return mediaInfo, fmt.Errorf("error reading collate path")
	}

	for _, file := range files {
		fileName := ParseFileName(mediaInfo, file.Name())
		mediaInfo.List[file.Name()] = fileName
	}

	return mediaInfo, nil
}
