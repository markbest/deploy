package utils

import (
	"log"
	"math"
	"os"
	"strconv"
	"sync"
)

var wg = &sync.WaitGroup{}

const chunkSize int64 = 2 << 20

// division zip file and upload
func ChunkFileUpload(zipPath string, prefix string, ssh *SSHClient) (chunkFiles []string, err error) {
	srcFile, err := os.Open(zipPath)
	if err != nil {
		return chunkFiles, err
	}
	defer srcFile.Close()

	fileInfo, err := os.Stat(zipPath)
	if err != nil {
		return chunkFiles, err
	}
	chunkNum := int(math.Ceil(float64(fileInfo.Size()) / float64(chunkSize)))
	log.Printf("压缩包分割成%d个小块文件", chunkNum)
	var i int64 = 1
	for ; i <= int64(chunkNum); i++ {
		wg.Add(1)
		b := make([]byte, chunkSize)
		srcFile.Seek((i-1)*(chunkSize), 0)
		if len(b) > int(fileInfo.Size()-(i-1)*chunkSize) {
			b = make([]byte, fileInfo.Size()-(i-1)*chunkSize)
		}
		srcFile.Read(b)

		chunkFile := prefix + "-" + strconv.FormatInt(i, 10) + ".zip"
		go uploadChunk(ssh, chunkFile, b, wg)
		chunkFiles = append(chunkFiles, chunkFile)
	}
	wg.Wait()
	return chunkFiles, nil
}

// upload chunk
func uploadChunk(ssh *SSHClient, chunkFile string, content []byte, wg *sync.WaitGroup) {
	dstFile, _ := ssh.Sftp.Create(chunkFile)
	dstFile.Write(content)
	log.Printf("小块文件%s上传完成", chunkFile)
	defer wg.Done()
	defer dstFile.Close()
}

// get merge chunk file command
func GetMergeFileCommand(chunkFiles []string, zipPath string) string {
	mergeFileCommand := "cat "
	if len(chunkFiles) > 0 {
		for _, f := range chunkFiles {
			mergeFileCommand = mergeFileCommand + f + " "
		}
		mergeFileCommand = mergeFileCommand + "> " + zipPath
	}
	return mergeFileCommand
}

// get delete chunk file command
func GetDeleteChunkFileCommand(chunkFiles []string, zipPath string) []string {
	removeFileCommand := make([]string, 0)
	if len(chunkFiles) > 0 {
		for _, f := range chunkFiles {
			removeFileCommand = append(removeFileCommand, "rm -f "+f)
		}
		removeFileCommand = append(removeFileCommand, "rm -f "+zipPath)
	}
	return removeFileCommand
}
