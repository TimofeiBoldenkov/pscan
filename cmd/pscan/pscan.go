package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/TimofeiBoldenkov/pscan/internal/tools"
)

func main() {
	portFlag := flag.String("port", "0-1023", "The target port")
	flag.StringVar(portFlag, "p", "0-1023", "The target port")
	showClosedPortsFlag := flag.Bool("show-closed", false, "Show closed ports")
	flag.BoolVar(showClosedPortsFlag, "c", false, "Show closed ports")
	protocol := flag.String("protocol", "tcp", "The protocol used in scanning (may be tcp, tcp4, tcp6, udp, udp4, udp6)")
	flag.StringVar(protocol, "t", "tcp", "The protocol used in scanning (may be tcp, tcp4, tcp6, udp, udp4, udp6)")
	flag.Parse()

	if !slices.Contains([]string{"tcp", "tcp4", "tcp6", "udp", "udp4", "udp6"}, *protocol) {
		fmt.Fprintf(os.Stderr, "Invalid protocol - %v\n", *protocol)
		os.Exit(1)
	}

	ports, err := tools.NewPorts(*portFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid ports - %v\n", *portFlag)
		os.Exit(1)
	}

	hosts := flag.Args()

	for port := range ports.Generate() {
		for _, host := range hosts {
			socket := net.JoinHostPort(host, strconv.FormatInt(int64(port), 10))
			timeout := 1 * time.Second

			_, err := net.DialTimeout(*protocol, socket, timeout)

			if os.IsTimeout(err) {
				fmt.Printf("%v: filtered\n", socket)
			} else if tools.IsRefused(err) {
				if *showClosedPortsFlag {
					fmt.Printf("%v: closed\n", socket)
				}
			} else {
				fmt.Printf("%v: open\n", socket)
			}
		}
	}
}
