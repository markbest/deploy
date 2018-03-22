package main

import (
	"bytes"
	. "github.com/markbest/deploy/conf"
	. "github.com/markbest/deploy/utils"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	outPut bytes.Buffer
	wg     = &sync.WaitGroup{}
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Parse config file
	if err := InitConfig(); err != nil {
		panic(err)
	}

	//Zip local files
	if len(Conf.Servers) > 0 {
		for _, server := range Conf.Servers {
			if len(server.Uploads) > 0 {
				for _, path := range server.Uploads {
					wg.Add(1)
					go handle(path, server, wg)
				}
			}
		}
	}
	wg.Wait()
}

func handle(path Uploads, server Server, wg *sync.WaitGroup) {
	//Create zip file
	timestamp := time.Now().Format("20060102030405")
	zipFile := timestamp + ".zip"
	zipPath := "./tmp/" + zipFile
	Zip(path.Local, zipPath, strings.Split(server.IgnoreDirs, ","))

	//Create ssh client
	ssh, err := NewSShClient(server.Host, server.Port, server.User, server.Password)
	if err != nil {
		panic(err)
	}

	//Upload file
	log.Printf("开始上传文件至服务器:%s", server.Host)
	chunkFiles, err := ChunkFileUpload(zipPath, path.Remote+timestamp, ssh)
	if err != nil {
		panic(err)
	}
	mergeFileCommand := GetMergeFileCommand(chunkFiles, path.Remote+zipFile)
	removeFileCommand := GetDeleteChunkFileCommand(chunkFiles, path.Remote+zipFile)

	//Run commands
	log.Printf("上传完毕开始解压文件")
	commands := make([]string, 0)
	commands = append(commands, mergeFileCommand)
	commands = append(commands, "/usr/bin/mkdir -p "+path.Remote+timestamp)
	commands = append(commands, "/usr/bin/unzip "+path.Remote+zipFile+" -d "+path.Remote+timestamp)
	ssh.Commands(commands, outPut)
	ssh.Commands(removeFileCommand, outPut)
	log.Printf("代码发布成功")
	defer wg.Done()
	defer ssh.Close()
}
