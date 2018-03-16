package main

import (
	"bytes"
	. "github.com/markbest/deploy/conf"
	. "github.com/markbest/deploy/utils"
	"os"
	"runtime"
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
	Zip(path.Local, zipPath)

	//Create ssh client
	ssh, err := NewSShClient(server.Host, server.Port, server.User, server.Password)
	if err != nil {
		panic(err)
	}

	//Upload file
	srcFile, err := os.Open(zipPath)
	if err != nil {
		panic(err)
	}

	dstFile, err := ssh.Sftp.Create(path.Remote + zipFile)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	//Run commands
	commands := make([]string, 0)
	commands = append(commands, "/usr/bin/mkdir -p "+path.Remote+timestamp)
	commands = append(commands, "/usr/bin/unzip "+path.Remote+zipFile+" -d "+path.Remote+timestamp)
	err = ssh.Commands(commands, outPut)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	defer ssh.Close()
}
