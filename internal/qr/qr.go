package qr

import (
	"fmt"
	"net"
	"strings"

	"rsc.io/qr"
)

func localIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no local IP found")
}

func buildURL(addr string) (string, error) {
	host, port := "", "8080"

	if strings.HasPrefix(addr, ":") {
		port = addr[1:]
	} else if strings.Contains(addr, ":") {
		parts := strings.SplitN(addr, ":", 2)
		host = parts[0]
		port = parts[1]
	} else {
		host = addr
	}

	if host == "" || host == "0.0.0.0" || host == "::" {
		ip, err := localIP()
		if err != nil {
			return "", err
		}
		host = ip
	}

	return fmt.Sprintf("http://%s:%s", host, port), nil
}

func printQR(url string) error {
	code, err := qr.Encode(url, qr.L)
	if err != nil {
		return err
	}

	chars := []string{" ", "▀", "▄", "█"}

	for y := 0; y < code.Size; y += 2 {
		for x := 0; x < code.Size; x++ {
			index := 0
			if code.Black(x, y) {
				index |= 1
			}
			if y+1 < code.Size && code.Black(x, y+1) {
				index |= 2
			}
			fmt.Print(chars[index])
		}
		fmt.Println()
	}
	return nil
}

func Print(addr string) {
	url, err := buildURL(addr)
	if err != nil {
		fmt.Printf("\nListening on %s\n\n", addr)
		return
	}

	if err := printQR(url); err != nil {
		fmt.Printf("\nListening on %s\n\n", url)
		return
	}

	fmt.Printf("%s\n", url)
}
