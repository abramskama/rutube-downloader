package rutube

import (
	"encoding/json"
	"fmt"
	"regexp"

	"rutube-downloader/requests"
)

type respOptions struct {
	VideoBalancer struct {
		M3U8 string `json:"m3u8"`
	} `json:"video_balancer"`
}

func GetVideoIDByURL(url string) (string, error) {
	r, err := regexp.Compile(`https://rutube.ru/video/(\w+)?/`)
	if err != nil {
		return "", err
	}
	finds := r.FindStringSubmatch(url)
	if len(finds) == 0 {
		return "", fmt.Errorf("invalid video url, url: %s", url)
	}
	videoId := finds[len(finds)-1]
	if videoId == "" {
		return "", fmt.Errorf("invalid video url, url: %s", url)
	}
	return videoId, nil
}

func GetPlaylistURL(videoID string) (string, error) {
	respBody, err := requests.GetResponseBody(fmt.Sprintf("https://rutube.ru/api/play/options/%s/?format=json", videoID))
	if err != nil {
		return "", err
	}

	var respJSON respOptions
	err = json.Unmarshal(respBody, &respJSON)
	if err != nil {
		return "", err
	}

	return respJSON.VideoBalancer.M3U8, nil
}
