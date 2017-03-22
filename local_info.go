package main

import (
	"net"
	"strings"
)

func getIP() string {
	ifAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	result := []string{}
	for _, ifAddr := range ifAddrs {
		addr := ifAddr.String()
		if i := strings.IndexByte(addr, '/'); i > 0 &&
			strings.IndexByte(addr, ':') == -1 && addr[:i] != `127.0.0.1` {
			result = append(result, addr[:i])
		}
	}
	return strings.Join(result, `,`)
}
