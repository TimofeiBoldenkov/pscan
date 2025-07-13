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

	"github.com/TimofeiBoldenkov/pscan/internal/scan"
	"github.com/TimofeiBoldenkov/pscan/internal/tools"
)

func main() {
	protocols := []string{"tcp", "tcp4", "tcp6"}

	portFlag := flag.String("port", "0-1023", "The target port")
	flag.StringVar(portFlag, "p", "0-1023", "The target port")
	timeoutFlag := flag.Float64("timeout", 20, "Timeout for a single connection in seconds")
	flag.Float64Var(timeoutFlag, "t", 20, "Timeout for a single connection in seconds")
	maxRequestsFlag := flag.Int("max-requests", 10, "The maximum amount of simultaneous requests to a single server")
	maxRoutinesFlag := flag.Int("max-routines", 10000, "The maxumum amount of goroutines")
	protocolHelpMessage := fmt.Sprintf("The protocol used in scanning (available protocols - %v)", protocols)
	protocolFlag := flag.String("protocol", "tcp", protocolHelpMessage)
	flag.Parse()

	*protocolFlag = strings.ToLower(*protocolFlag)
	if !slices.Contains(protocols, *protocolFlag) {
		fmt.Fprintf(os.Stderr, "Invalid protocol - %v\n", *protocolFlag)
		os.Exit(1)
	}

	ports, err := tools.NewPorts(*portFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid ports - %v\n", *portFlag)
		os.Exit(1)
	}

	hosts := flag.Args()

	ch := make(chan struct{}, tools.Min(*maxRequestsFlag*len(hosts), *maxRoutinesFlag))

	type portInfo struct {
		port  int
		state scan.PortState
	}
	result := make(map[string][]portInfo, len(hosts)*3/2)

	fmt.Println("Scanning...")
	fmt.Println()

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for port := range ports.Generate() {
		for _, host := range hosts {
			ch <- struct{}{}
			wg.Add(1)
			go func(host string, port int, protocol string, timeout time.Duration) {
				defer func() {
					<-ch
					wg.Done()
				}()
				portState, err := scan.ScanHostPort(host, strconv.FormatInt(int64(port), 10), protocol, timeout)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				} else if portState != scan.Closed {
					mutex.Lock()
					result[host] = append(result[host], portInfo{port, portState})
					mutex.Unlock()
				}
			}(host, port, *protocolFlag, time.Duration(*timeoutFlag*float64(time.Second)))
		}
	}

	wg.Wait()
	close(ch)

	for host, infos := range result {
		fmt.Printf("%v:\n", host)
		for _, info := range infos {
			fmt.Printf("\t%v\t%v\n", info.port, info.state)
		}
		fmt.Println()
	}
}
