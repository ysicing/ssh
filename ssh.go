// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package ssh

import (
	"os"

	"github.com/ysicing/logger"
)

func (ss *SSH) Run(host string, cmd string) {
	session, err := ss.Connect(host)
	defer func() {
		if r := recover(); r != nil {
			logger.Slog.Errorf("%v create ssh session failed, %v", host, err)
		}
	}()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if err := session.Run(cmd); err != nil {
		logger.Slog.Errorf("%v exec cmd(%v) failed, %v", host, cmd, err)
		os.Exit(0)
	}
}
