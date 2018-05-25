package utils

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type File struct {
	Name string
	Body string
}

// Get file list
func getFileList(localPath string, ignoreDirs []string) (files []string, err error) {
	PthSep := string(os.PathSeparator)
	err = filepath.Walk(localPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		flag := judgeIgnore(localPath, path, ignoreDirs)
		if !flag {
			if f.IsDir() {
				files = append(files, path+PthSep)
			} else {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return files, err
	}
	return files, err
}

// Judge file is ignore or not
func judgeIgnore(path string, file string, ignoreDirs []string) bool {
	flag := false
	if len(ignoreDirs) > 0 {
		for _, d := range ignoreDirs {
			if strings.HasPrefix(file, path+d) {
				flag = true
				break
			}
		}
	}
	return flag
}

// Zip compress
func Zip(localPath string, zipPath string, ignoreDirs []string) {
	log.Printf("开始打包目录:%s", localPath)
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	zipFiles := make([]File, 0)
	files, err := getFileList(localPath, ignoreDirs)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		content, _ := ioutil.ReadFile(file)
		zipFiles = append(zipFiles, File{file, string(content)})
	}

	for _, file := range zipFiles {
		f, err := w.CreateHeader(
			&zip.FileHeader{
				Name:     strings.Replace(file.Name, localPath, "", 1),
				Modified: time.Now(),
			},
		)
		if err != nil {
			panic(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			panic(err)
		}
	}

	err = w.Close()
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(zipPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	buf.WriteTo(f)
}
