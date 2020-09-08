// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package ssh

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/ysicing/go-utils/exfile"
	"github.com/ysicing/logger"
	"golang.org/x/crypto/ssh"
)

func (ss *SSH) sshAuth(passwd, pkfile string) (auth []ssh.AuthMethod) {
	if strings.Contains(pkfile, "~") {
		home, _ := homedir.Dir()
		pkfile = strings.ReplaceAll(pkfile, "~", home)
	}
	if exfile.CheckFileExistsv2(pkfile) {
		pkfiledata, err := ioutil.ReadFile(pkfile)
		if err != nil {
			logger.Exitf("readv pkfile err: %v", err)
		}
		pk, err := ssh.ParsePrivateKey(pkfiledata)
		if err == nil {
			auth = append(auth, ssh.PublicKeys(pk))
		}
	}
	if len(passwd) > 0 {
		auth = append(auth, ssh.Password(passwd))
	}
	return auth
}

func (ss *SSH) addrReformat(host string) string {
	if strings.Index(host, ":") == -1 {
		host = fmt.Sprintf("%s:22", host)
	}
	return host
}

func (ss *SSH) connect(host string) (*ssh.Client, error) {
	auth := ss.sshAuth(ss.Password, ss.PkFile)
	config := ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	DefaultTimeout := time.Duration(1) * time.Minute
	if ss.Timeout == nil {
		ss.Timeout = &DefaultTimeout
	}
	clientConfig := &ssh.ClientConfig{
		User:    ss.User,
		Auth:    auth,
		Timeout: *ss.Timeout,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr := ss.addrReformat(host)
	return ssh.Dial("tcp", addr, clientConfig)
}

func (ss *SSH) Connect(host string) (*ssh.Session, error) {
	client, err := ss.connect(host)
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}
	return session, nil
}
