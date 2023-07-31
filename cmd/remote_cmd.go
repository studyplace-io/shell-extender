package cmd

import (
	"fmt"
	"github.com/practice/shell_extender/pkg/remote_command"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func remoteExecShellCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remoteExec",
		Short: "exec shell script for remote server",
		Long:  "exec shell script for remote server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &config{
				host: host,
				user: user,
				password: password,
				port: port,
			}
			err := remote_command.RunRemoteNode(cfg.user, cfg.password, cfg.host, cfg.port, readFile(script))
			return err
		},
	}
	return cmd
}

func readFile(path string) string {
	// 读取文件内容
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("read file error：", err)
		return ""
	}

	return string(data)
}