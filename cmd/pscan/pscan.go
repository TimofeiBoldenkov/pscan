package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/TimofeiBoldenkov/pscan/internal/tools"
	"github.com/TimofeiBoldenkov/pscan/internal/scan"
)

func main() {
	portFlag := flag.String("port", "0-1023", "The target port")
	flag.StringVar(portFlag, "p", "0-1023", "The target port")
	showClosedPortsFlag := flag.Bool("show-closed", false, "Show closed ports")
	flag.BoolVar(showClosedPortsFlag, "c", false, "Show closed ports")
	protocolFlag := flag.String("protocol", "tcp", "The protocol used in scanning (may be tcp, tcp4, tcp6, udp, udp4, udp6)")
	flag.StringVar(protocolFlag, "t", "tcp", "The protocol used in scanning (may be tcp, tcp4, tcp6, udp, udp4, udp6)")
	maxRequestsFlag := flag.Int("max-requests", 10, "The maximum semultaneous requests to a server")
	flag.Parse()

	*protocolFlag = strings.ToLower(*protocolFlag)
	if !slices.Contains([]string{"tcp", "tcp4", "tcp6", "udp", "udp4", "udp6"}, *protocolFlag) {
		fmt.Fprintf(os.Stderr, "Invalid protocol - %v\n", *protocolFlag)
		os.Exit(1)
	}

	ports, err := tools.NewPorts(*portFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid ports - %v\n", *portFlag)
		os.Exit(1)
	}

	hosts := flag.Args()

	var wg sync.WaitGroup

	const MAX_ROUTINES int = 10000
	ch := make(chan struct{}, tools.Min(*maxRequestsFlag * len(hosts), MAX_ROUTINES))

	for port := range ports.Generate() {
		for _, host := range hosts {
			wg.Add(1)
			ch <- struct{}{}
			go func() {
				defer func() {
					<-ch
					wg.Done() 
				}()
				scan.ScanSocket(host, strconv.FormatInt(int64(port), 10), *protocolFlag, 20 * time.Second, *showClosedPortsFlag)
			}()
		}
	}

	wg.Wait()
	close(ch)
}
