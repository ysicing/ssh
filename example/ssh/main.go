// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package main

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/go-utils/extime"
	"github.com/ysicing/logger"
	"github.com/ysicing/ssh"
)

var ips []string
var user string
var pass string
var pkfile string
var xcmd string

var rootCmd = &cobra.Command{
	Use:   "ssh",
	Short: "命令行工具",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(ips) < 1 {
			logger.Exit("ip不允许为空")
		}
		if len(pass) == 0 && len(pkfile) == 0 {
			logger.Exit("认证信息为空, 请指定密码或者私钥")
		}
		if len(xcmd) < 1 {
			logger.Slog.Debug("执行命令为空")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, ip := range ips {
			st := extime.NowUnix()
			logger.Slog.Infof("=== [ %v ] 开始执行命令 ===", ip)
			s := ssh.SSH{
				User:     user,
				Password: pass,
				PkFile:   pkfile,
			}
			s.Run(ip, xcmd)
			et := extime.NowUnix()
			logger.Slog.Infof("=== [ %v ] 完成执行命令, 耗时: %v ===\n", ip, et-st)
		}
	},
}

func init() {
	cfg := &logger.LogConfig{Simple: true}
	logger.InitLogger(cfg)
	rootCmd.PersistentFlags().StringArrayVarP(&ips, "ip", "", nil, "ip地址eg: 127.0.0.1, 127.0.0.1:2222")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "root", "用户")
	rootCmd.PersistentFlags().StringVarP(&pass, "pass", "p", "", "密码")
	rootCmd.PersistentFlags().StringVarP(&pkfile, "pkfile", "k", "", "私钥")
	rootCmd.PersistentFlags().StringVarP(&xcmd, "c", "", "", "命令")
}

func main() {
	rootCmd.Execute()
}
