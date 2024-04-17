package m3u8

import (
	"errors"
	"log"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/samber/lo"
	ffmpeg "github.com/u2takey/ffmpeg-go"

	"rutube-downloader/requests"
)

func GetSegmentsURLs(playListUrl string) ([]string, error) {
	resp, err := requests.GetResponse(playListUrl)
	if err != nil {
		// log.Printf("Error, cant get playlist, HTTP_CODE = %d", resp.StatusCode)
		return nil, err
	}

	p, listType, err := m3u8.DecodeFrom(resp.Body, false)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()
	var maxResolutionVariant *m3u8.Variant
	switch listType {
	case m3u8.MEDIA:
		//maxResolutionVariant = p.(*m3u8.MediaPlaylist)
		return nil, errors.New("")
	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		maxResolutionVariant = lo.MaxBy(masterpl.Variants, func(p1, p2 *m3u8.Variant) bool {
			return p1.Resolution > p2.Resolution
		})
	}

	// Get segments files
	resp, err = requests.GetResponse(maxResolutionVariant.URI)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()
	p, listType, err = m3u8.DecodeFrom(resp.Body, false)
	if err != nil {
		return nil, err
	}
	var urls []string
	switch listType {
	case m3u8.MEDIA:
		playlist := p.(*m3u8.MediaPlaylist)
		log.Println(playlist)
		splitted := strings.Split(maxResolutionVariant.URI, "/")
		splitted = splitted[:len(splitted)-1]
		baseUri := strings.Join(splitted, "/")
		segments := playlist.GetAllSegments()
		urls = lo.Map(segments, func(item *m3u8.MediaSegment, idx int) string {
			return baseUri + "/" + item.URI
		})
	case m3u8.MASTER:
		return nil, errors.New("")
	}

	return urls, nil
}

func JoinSegments(srcPaths []string, dstPath string) error {
	joined := "concat:" + strings.Join(srcPaths, "|")
	return ffmpeg.Input(joined).
		Output(dstPath, ffmpeg.KwArgs{"acodec": "copy", "vcodec": "copy"}).
		OverWriteOutput().ErrorToStdOut().Run()
}
