package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

/*
=== Утилита wget ===
Реализовать утилиту wget с возможностью скачивать сайты целиком
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type file struct {
	io.Reader
	name    string
	size    int64
	counter int64
}

func main() {
	wg := &sync.WaitGroup{}
	var fileList []*file

	for _, uRL := range os.Args[1:] {
		link, bodySize, err := parseURL(uRL)
		if err == nil {
			file := &file{size: bodySize}
			fileList = append(fileList, file)
			wg.Add(1)
			go file.download(link, wg)
		} else {
			log.Printf("wget: error parsing url %v: %v\n", uRL, err)
		}
	}

	endCh := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go print(ctx, fileList, endCh)

	wg.Wait()
	cancel()
	<-endCh
}

func (file *file) download(link string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(link)
	if err != nil {
		log.Printf("Err of connecting to url :%v\n", err)
		return
	}
	defer resp.Body.Close()

	fileForSave := path.Base(resp.Request.URL.Path)

	f, err := os.Create(fileForSave)
	if err != nil {
		log.Printf("wget: error creating file :%v\n", err)
		return
	}

	file.name = fileForSave
	file.Reader = resp.Body

	_, err = io.Copy(f, file)
	if err != nil {
		log.Printf("Error of reading from %s\n", link)
		return
	}
}

func print(ctx context.Context, fileList []*file, endCh chan struct{}) {
	defer close(endCh)

	print := func() {
		fmt.Print("\r")
		for _, file := range fileList {
			if file.size != 0 {
				percentage := float64(file.counter) / float64(file.size) * 100
				fmt.Printf("%.2f %% ", percentage)
			} else {
				fmt.Print("0.00 %% ")
			}
		}
	}

	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			print()
			return
		case <-t.C:
			print()
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func parseURL(link string) (parsedLink string, bodySize int64, err error) {
	_, err = url.ParseRequestURI(link)
	if err != nil {
		return
	}
	resp, err := http.Get(link)
	if err != nil {
		err = fmt.Errorf("wget: error connecting url :%v", err)
		return
	}
	defer resp.Body.Close()

	bodySize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 0)
	if err != nil {
		err = fmt.Errorf("wget: content length = %v :%v", bodySize, err)
		return
	}

	parsedLink = link
	return
}

func (file *file) Read(p []byte) (int, error) {
	n, err := file.Reader.Read(p)
	if err == nil {
		file.counter += int64(n)
	}

	return n, err
}
