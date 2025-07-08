package tools

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type portsEntry struct {
	firstPort int
	lastPort  int
}

func newPortsEntry(portsEntryStr string) (portsEntry, bool) {
	separator := strings.Index(portsEntryStr, "-")
	if separator == -1 {
		firstPort, err := strconv.Atoi(portsEntryStr)
		if err != nil {
			return portsEntry{}, false
		} else {
			return portsEntry{firstPort, firstPort}, true
		}
	} else {
		firstPort, firstErr := strconv.Atoi(portsEntryStr[0:separator])
		lastPort, lastErr := strconv.Atoi(portsEntryStr[separator+1:])
		if firstErr != nil || lastErr != nil || firstPort > lastPort ||
			firstPort < 0 || firstPort > 65535 || lastPort < 0 || lastPort > 65535 {
			return portsEntry{}, false
		} else {
			return portsEntry{firstPort, lastPort}, true
		}
	}
}

type Ports struct {
	ports []portsEntry
	value int
}

func NewPorts(portsStr string) (Ports, error) {
	var retPorts Ports

	entryBegin := 0
	for {
		entryEnd := strings.Index(portsStr[entryBegin:], ",")
		if entryEnd == -1 {
			entryEnd = len(portsStr)
		} else {
			entryEnd += entryBegin
		}

		entry, ok := newPortsEntry(portsStr[entryBegin:entryEnd])
		if !ok {
			errMessage := fmt.Sprintf("Invalid entry - %v", portsStr[entryBegin:entryEnd])
			return Ports{}, errors.New(errMessage)
		} else {
			retPorts.ports = append(retPorts.ports, entry)
			if entryEnd == len(portsStr) {
				return retPorts, nil
			}
			entryBegin = entryEnd + 1
		}
	}
}

func (ports *Ports) Next() bool {
	if len(ports.ports) == 0 {
		return false
	} else {
		ports.value = ports.ports[0].firstPort
		if ports.ports[0].firstPort == ports.ports[0].lastPort {
			ports.ports = ports.ports[1:]
		} else {
			ports.ports[0].firstPort++
		}
		return true
	}
}

func (ports *Ports) Value() int {
	return ports.value
}

func (ports *Ports) Generate() <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for ports.Next() {
			ch <- ports.Value()
		}
	}()

	return ch
}

func (ports *Ports) Copy() Ports {
	dataCopy := make([]portsEntry, len(ports.ports))
	copy(dataCopy, ports.ports)
	return Ports{dataCopy, ports.value}
}
