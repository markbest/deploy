package utils

import (
	"log"
	"math"
	"os"
	"strconv"
)

const chunkSize int64 = 4 << 20

//division zip file and upload
func ChunkFileUpload(zipPath string, prefix string, ssh *SSHClient) (chunkFiles []string, err error) {
	srcFile, err := os.Open(zipPath)
	if err != nil {
		return chunkFiles, err
	}

	fileInfo, err := os.Stat(zipPath)
	if err != nil {
		return chunkFiles, err
	}
	chunkNum := int(math.Ceil(float64(fileInfo.Size()) / float64(chunkSize)))
	log.Printf("压缩包分割成%d个小块文件", chunkNum)

	b := make([]byte, chunkSize)
	var i int64 = 1
	for ; i <= int64(chunkNum); i++ {
		srcFile.Seek((i-1)*(chunkSize), 0)
		if len(b) > int((fileInfo.Size() - (i-1)*chunkSize)) {
			b = make([]byte, fileInfo.Size()-(i-1)*chunkSize)
		}
		srcFile.Read(b)

		chunkFile := prefix + "-" + strconv.FormatInt(i, 10) + ".zip"
		dstFile, err := ssh.Sftp.Create(chunkFile)
		if err != nil {
			panic(err)
		}
		chunkFiles = append(chunkFiles, chunkFile)

		dstFile.Write(b)
		dstFile.Close()
		log.Printf("小块文件%s上传完成", chunkFile)
	}
	srcFile.Close()
	return chunkFiles, nil
}

//get merge chunk file command
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

//get delete chunk file command
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
