package utils

import (
	"bytes"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"strconv"
)

type SSHClient struct {
	Client *ssh.Client
	Sftp   *sftp.Client
}

//New ssh client connect
func NewSShClient(ip string, port int64, user string, password string) (sc *SSHClient, err error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", ip+":"+strconv.FormatInt(port, 10), config)
	if err != nil {
		return sc, err
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return sc, err
	}

	sc = &SSHClient{Client: client, Sftp: sftpClient}
	return sc, nil
}

//Run commands
func (sc *SSHClient) Commands(commands []string, output bytes.Buffer) error {
	if len(commands) > 0 {
		for _, command := range commands {
			session, err := sc.Client.NewSession()
			if err != nil {
				return err
			}
			defer session.Close()
			session.Stdout = &output
			err = session.Run(command)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//Close ssh connect
func (sc *SSHClient) Close() {
	sc.Sftp.Close()
}
