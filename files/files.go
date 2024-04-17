package files

import (
	"io"
	"log"
	"os"
	"sync"

	"rutube-downloader/requests"
)

type FileOption struct {
	Url  string
	Path string
}

func BulkDownload(parallelLimit int, files []FileOption) error {
	ch := make(chan struct{}, parallelLimit)
	wg := sync.WaitGroup{}
	for _, file := range files {
		wg.Add(1)
		go func(file FileOption) {
			ch <- struct{}{}
			err := download(file)
			if err != nil {
				panic(err)
			}
			wg.Done()
			<-ch
		}(file)
	}
	wg.Wait()
	return nil
}

func download(file FileOption) error {
	log.Printf("%s - START downloads file", file.Path)

	resp, err := requests.GetResponse(file.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(file.Path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	log.Printf("%s - FINISH downloads video", file.Path)
	return nil
}
