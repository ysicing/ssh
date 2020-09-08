// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package ssh

import "time"

type SSH struct {
	User string
	Password string
	PkFile string
	Timeout *time.Duration
}