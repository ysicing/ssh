// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package ssh

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/ysicing/ext/utils/exmisc"

	"github.com/ysicing/ext/logger"
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

func readPipe(host string, pipe io.Reader, isErr bool, rundebug ...bool) {
	r := bufio.NewReader(pipe)
	for {
		line, _, err := r.ReadLine()
		if line == nil {
			return
		} else if err != nil {
			logger.Slog.Infof("[%s] %s", exmisc.SGreen(host), line)
			logger.Slog.Errorf("[ssh] [%s] %s", exmisc.SRed(host), err)
			return
		} else {
			if isErr {
				if len(rundebug) > 0 && rundebug[0] {
					logger.Slog.Errorf("[%s] %s", exmisc.SRed(host), line)
				} else {
					msg, _ := fmt.Printf("%s", line)
					fmt.Println(msg)
				}
			} else {
				if len(rundebug) > 0 && rundebug[0] {
					logger.Slog.Infof("[%s] %s", exmisc.SGreen(host), line)
				} else {
					msg, _ := fmt.Printf("%s", line)
					fmt.Println(msg)
				}
			}
		}
	}
}

func (ss *SSH) CmdAsync(host string, cmd string, wg *sync.WaitGroup, rundebug ...bool) error {
	defer wg.Done()
	fmt.Printf("[ssh][%s] ➜   %s\n", exmisc.SGreen(host), cmd)
	session, err := ss.Connect(host)
	if err != nil {
		logger.Slog.Errorf("[ssh][%s]Error create ssh session failed,%s", host, err)
		return err
	}
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		logger.Slog.Errorf("[ssh][%s]Unable to request StdoutPipe(): %s", host, err)
		return err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		logger.Slog.Errorf("[ssh][%s]Unable to request StderrPipe(): %s", host, err)
		return err
	}
	if err := session.Start(cmd); err != nil {
		logger.Slog.Errorf("[ssh][%s]Unable to execute command: %s", host, err)
		return err
	}
	doneout := make(chan bool, 1)
	doneerr := make(chan bool, 1)
	go func() {
		readPipe(host, stderr, true, rundebug...)
		doneerr <- true
	}()
	go func() {
		readPipe(host, stdout, false, rundebug...)
		doneout <- true
	}()
	<-doneerr
	<-doneout
	return session.Wait()
}
