// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package main

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/utils/extime"
	"github.com/ysicing/ssh"
	"sync"
)

var (
	ips    []string
	user   string
	pass   string
	pkfile string
	xcmd   string
	mp     bool // Multiple processes 多进程
)

var rootCmd = &cobra.Command{
	Use:   "ssh",
	Short: "命令行工具",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(ips) < 1 {
			logger.Slog.Exit0("ip不允许为空")
		}
		if len(pass) == 0 && len(pkfile) == 0 {
			logger.Slog.Exit0("认证信息为空, 请指定密码或者私钥")
		}
		if len(xcmd) < 1 {
			logger.Slog.Debug("执行命令为空")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		s := ssh.SSH{
			User:     user,
			Password: pass,
			PkFile:   pkfile,
		}
		if mp {
			var wg sync.WaitGroup
			for _, ip := range ips {
				wg.Add(1)
				s.CmdAsync(ip, xcmd, &wg)
			}
			wg.Wait()
		} else {
			for _, ip := range ips {
				st := extime.NowUnix()
				logger.Slog.Infof("=== [ %v ] 开始执行命令 ===", ip)
				s.Run(ip, xcmd)
				et := extime.NowUnix()
				logger.Slog.Infof("=== [ %v ] 完成执行命令, 耗时: %v ===\n", ip, et-st)
			}
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
	rootCmd.PersistentFlags().BoolVarP(&mp, "mp", "", true, "多进程")
}

func main() {
	rootCmd.Execute()
}
