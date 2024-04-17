package main

import (
	"fmt"
	"log"

	"github.com/samber/lo"

	"rutube-downloader/files"
	"rutube-downloader/m3u8"
	"rutube-downloader/rutube"
)

const (
	VideoUrl             = "https://rutube.ru/video/fa15cefeb832b515d6c283e2a2ba87c8/"
	DirectoryForDownload = "./downloads"
	VideoPartPath        = DirectoryForDownload + "/%s_part_%d.ts"
	ResultVideoPath      = DirectoryForDownload + "/%s_result_video.mp4"

	ParallelDownloadsLimit = 10
)

func main() {
	log.Printf("Get video id, url: %s", VideoUrl)
	videoID, err := rutube.GetVideoIDByURL(VideoUrl)
	if err != nil {
		panic(err)
	}

	log.Printf("Get video playlist url (m3u8), id: %s", videoID)
	playListURL, err := rutube.GetPlaylistURL(videoID)
	if err != nil {
		panic(err)
	}

	log.Printf("Get video segments urls, playListURL: %s", playListURL)
	urls, err := m3u8.GetSegmentsURLs(playListURL)
	if err != nil {
		panic(err)
	}

	log.Printf("Download segments, cnt: %d", len(urls))
	filesOptions := lo.Map(urls, func(url string, idx int) files.FileOption {
		return files.FileOption{
			Url:  url,
			Path: fmt.Sprintf(VideoPartPath, videoID, idx),
		}
	})
	err = files.BulkDownload(ParallelDownloadsLimit, filesOptions)
	if err != nil {
		panic(err)
	}

	log.Printf("Join segments")
	paths := lo.Map(filesOptions, func(file files.FileOption, _ int) string {
		return file.Path
	})
	err = m3u8.JoinSegments(paths, fmt.Sprintf(ResultVideoPath, videoID))
	if err != nil {
		panic(err)
	}

	log.Printf("Video downloaded, path: %s", fmt.Sprintf(ResultVideoPath, videoID))
}
