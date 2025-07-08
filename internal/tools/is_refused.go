package tools

import "net"

func IsRefused(err error) bool {
	_, ok := err.(*net.OpError)
	return ok
}
