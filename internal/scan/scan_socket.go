package scan

import (
	"github.com/TimofeiBoldenkov/pscan/internal/tools"
	"time"
	"net"
	"os"
	"fmt"
)

func ScanSocket(host string, port string, protocol string, timeout time.Duration, showClosed bool) {
	socket := net.JoinHostPort(host, port)

	_, err := net.DialTimeout(protocol, socket, timeout)

	if os.IsTimeout(err) {
		fmt.Printf("%v: filtered\n", socket)
	} else if tools.IsRefused(err) {
		if showClosed {
			fmt.Printf("%v: closed\n", socket)
		}
	} else if err == nil {
		fmt.Printf("%v: open\n", socket)
	} else {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
	}
}
