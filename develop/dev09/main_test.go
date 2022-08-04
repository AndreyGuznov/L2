package main

import (
	"fmt"
	"os"
	"path"
	"sync"
	"testing"
)

func TestDownload(t *testing.T) {

	var tests = []struct {
		inpURL string
		result bool
	}{
		{"https://wordpress.org/latest.zip", true},
		{"qwerty.com", false},
		{"https://google.com/goog.txt", true},
		{"https://yandex", false},
	}

	wg := &sync.WaitGroup{}
	var fileList []*file

	for _, test := range tests {
		link, bodySize, err := parseURL(test.inpURL)
		if err == nil {
			file := &file{size: bodySize}
			fileList = append(fileList, file)
			wg.Add(1)
			go file.download(link, wg)
		} else {
			t.Logf("Error while parsing url %v\n", test.inpURL)
		}
	}
	wg.Wait()

	for _, test := range tests {
		uRL := test.inpURL
		link := path.Base(uRL)

		var maxSize int64
		for i := range fileList {
			if fileList[i].name == link {
				maxSize = fileList[i].size
				break
			}
		}

		forRes, err := fileDownloaded(link, maxSize)
		if forRes != test.result {
			t.Errorf("Download(%q) = %v : %v", uRL, forRes, err.Error())
		}
	}
}

func fileDownloaded(fileName string, wantedSize int64) (bool, error) {
	file, err := os.Stat(fileName)
	if err != nil {
		return false, fmt.Errorf("%s: %v", fileName, err)
	}

	if file.Size() != wantedSize {
		return false, fmt.Errorf("%s err of downloading", fileName)
	}

	return true, nil
}
