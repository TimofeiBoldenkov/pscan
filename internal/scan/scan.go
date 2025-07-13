package scan

import (
	"errors"
	"net"
	"os"
	"time"
	"fmt"

	"github.com/TimofeiBoldenkov/pscan/internal/tools"
)

type PortState uint8

const (
	Open PortState = iota
	Closed
	Filtered
	OpenFiltered
	Unknown
)

func (ps PortState) String() string {
	switch ps {
	case Open:
		return "open"
	case Closed:
		return "closed"
	case Filtered:
		return "filtered"
	case OpenFiltered:
		return "OpenFiltered"
	default:
		return "unknown"
	}
}

func ScanHostPort(host string, port string, protocol string, timeout time.Duration) (PortState, error) {
	socket := net.JoinHostPort(host, port)

	if protocol == "tcp" || protocol == "tcp4" || protocol == "tcp6" {
		_, err := net.DialTimeout(protocol, socket, timeout)
		if err == nil {
			return Open, nil
		} else if os.IsTimeout(err) {
			return Filtered, nil
		} else if tools.IsRefused(err) {
			return Closed, nil
		} else {
			return Unknown, err
		}
	} else {
		errMessage := fmt.Sprintf("unsupported protocol - %v", protocol)
		return Unknown, errors.New(errMessage)
	}
}
