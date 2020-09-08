// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package ssh

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	*SSH
	SSHClient *ssh.Client
	SSHSession *ssh.Session
	SFTPClient *sftp.Client
}