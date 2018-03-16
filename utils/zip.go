package utils

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name string
	Body string
}

//Get file list
func getFileList(path string) (files []string, err error) {
	PthSep := string(os.PathSeparator)
	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			files = append(files, path+PthSep)
		} else {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return files, err
	}
	return files, err
}

//Zip compress
func Zip(path string, name string) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	zipFiles := make([]File, 0)
	files, err := getFileList(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		content, _ := ioutil.ReadFile(file)
		zipFiles = append(zipFiles, File{file, string(content)})
	}

	for _, file := range zipFiles {
		f, err := w.Create(strings.Replace(file.Name, path, "", 1))
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

	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	buf.WriteTo(f)
}
