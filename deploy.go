package main

import (
	"bytes"
	"log"
	"runtime"
	"strings"
	"time"

	. "github.com/markbest/deploy/conf"
	. "github.com/markbest/deploy/utils"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse config file
	if err := InitConfig(); err != nil {
		panic(err)
	}

	// Zip local files
	if len(Conf.Servers) > 0 {
		for _, server := range Conf.Servers {
			if len(server.Uploads) > 0 {
				for _, path := range server.Uploads {
					handleUpload(path, server)
				}
			}
		}
	}
}

func handleUpload(path Uploads, server Server) {
	var outPut bytes.Buffer

	// Create zip file
	timestamp := time.Now().Format("20060102030405")
	zipFile := timestamp + ".zip"
	zipPath := "./tmp/" + zipFile
	Zip(path.Local, zipPath, strings.Split(server.IgnoreDirs, ","))

	// Create ssh client
	ssh, err := NewSShClient(server.Host, server.Port, server.User, server.Password)
	if err != nil {
		panic(err)
	}
	defer ssh.Close()

	// Run pre commands
	if server.PreCommands != "" {
		preCommands := strings.Split(server.PreCommands, ",")
		ssh.Commands(preCommands, outPut)
	}

	// Upload file
	log.Printf("开始上传文件至服务器:%s", server.Host)
	chunkFiles, err := ChunkFileUpload(zipPath, path.Remote+timestamp, ssh)
	if err != nil {
		panic(err)
	}
	mergeFileCommand := GetMergeFileCommand(chunkFiles, path.Remote+zipFile)
	removeFileCommand := GetDeleteChunkFileCommand(chunkFiles, path.Remote+zipFile)

	// Unzip file
	log.Printf("上传完毕开始解压文件")
	commands := make([]string, 0)
	commands = append(commands, mergeFileCommand)
	commands = append(commands, "/usr/bin/mkdir -p "+path.Remote+timestamp)
	commands = append(commands, "/usr/bin/unzip "+path.Remote+zipFile+" -d "+path.Remote+timestamp)
	ssh.Commands(commands, outPut)
	ssh.Commands(removeFileCommand, outPut)

	// Run post commands
	if server.PreCommands != "" {
		postCommands := strings.Split(server.PreCommands, ",")
		ssh.Commands(postCommands, outPut)
	}

	log.Printf("代码发布成功")
}
