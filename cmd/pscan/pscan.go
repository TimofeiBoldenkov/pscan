package main

import (
	"fmt"
	"flag"

	"github.com/TimofeiBoldenkov/pscan/internal/tools"
)

func main() {
	portFlag := flag.String("port", "0-1023", "The target port")
	flag.StringVar(portFlag, "p", "0-1023", "The target port")
	flag.Parse()

	ports, err := tools.NewPorts(*portFlag)
	if err != nil {
		fmt.Println(err.Error())
	}

	for port := range ports.Generate() {
		fmt.Println(port)
	}
}
