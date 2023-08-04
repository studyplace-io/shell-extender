package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "opentelemetry-test-server",
	Long:  "",
}

var (
	host     string
	user     string
	password string
	port     int
	script   string
)

func init() {
	runCmd.PersistentFlags().StringVarP(&host, "host", "i", "", "host ip")
	runCmd.PersistentFlags().StringVarP(&user, "user", "u", "root", "remote user name")
	runCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "remote user password")
	runCmd.PersistentFlags().IntVarP(&port, "port", "P", 22, "remote user password")
	runCmd.PersistentFlags().StringVarP(&script, "script", "s", "", "bash shell script")
	runCmd.AddCommand(remoteExecShellCmd(), remoteExecPingCmd(), remoteCommandLineCmd())
}

func Execute() {
	if err := runCmd.Execute(); err != nil {
		fmt.Printf("cmd err: %s\n", err)
		os.Exit(1)
	}
}
