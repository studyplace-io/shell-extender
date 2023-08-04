package cmd

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

func remoteExecPingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remotePing",
		Short: "ping for remote server",
		Long:  "ping for remote server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &config{
				host:     host,
				user:     user,
				password: password,
				port:     port,
			}
			pinger, err := ping.NewPinger(cfg.host)
			if err != nil {
				panic(err)
			}
			//if runtime.GOOS=="windows"{
			//	pinger.SetPrivileged(true )
			//}
			defer pinger.Stop()

			pinger.OnRecv = func(pkt *ping.Packet) {
				fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
					pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
			}

			pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
				fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
					pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
			}

			pinger.OnFinish = func(stats *ping.Statistics) {
				fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
				fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
					stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
				fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
					stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
			}

			fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

			pinger.Count = 5 // 最多五次
			err = pinger.Run()
			if err != nil {
				panic(err)
			}
			return err
		},
	}
	return cmd
}
