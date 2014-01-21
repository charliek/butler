package iptables

import (
	"fmt"
	"github.com/charliek/butler/task"
	"strings"
)

type PortLocation struct {
	Chain string
	Index int
}

type IPAddress struct {
	Interface string
	Address   string
}

var host_ip = "192.168.70.1"

// TODO this should be pulled dynamically
var proxy_ips = []string{"10.0.2.15", "192.168.70.99", "127.0.0.1"}

func parseNatRules(natOut string, port int) []*PortLocation {
	portLocations := make([]*PortLocation, 0, 5)
	currentChain := ""
	currentIndex := 0

	for _, line := range strings.Split(natOut, "\n") {
		if strings.HasPrefix(line, "Chain") {
			currentChain = strings.Split(line, " ")[1]
			currentIndex = 0
		} else if strings.HasPrefix(line, "DNAT") || strings.HasPrefix(line, "SNAT") {
			currentIndex += 1
			if strings.Contains(line, fmt.Sprintf("to:%s:%d", host_ip, port)) {
				portLocations = append(portLocations, &PortLocation{currentChain, currentIndex})
			}
		}
	}

	return portLocations
}

func parseIsMasquerade(natOut string) bool {
	currentChain := ""
	for _, line := range strings.Split(natOut, "\n") {
		if strings.HasPrefix(line, "Chain") {
			currentChain = strings.Split(line, " ")[1]
		} else if currentChain == "POSTROUTING" && strings.HasPrefix(line, "MASQUERADE") {
			return true
		}
	}
	return false
}

func getNatRules() (string, error) {
	return task.ExecuteStringTask("sudo iptables -t nat -L")
}

func IsLocal(port int) (bool, error) {
	natOut, err := getNatRules()
	if err != nil {
		return false, err
	}

	// Exit if we already have iptable rules setup for this port
	portLocations := parseNatRules(natOut, port)
	return len(portLocations) != 0, nil
}

func RunVagrant(port int) error {
	natOut, err := getNatRules()
	if err != nil {
		return err
	}

	// Exit if we already have iptable rules setup for this port
	portLocations := parseNatRules(natOut, port)
	if len(portLocations) == 0 {
		return nil
	}

	for i := len(portLocations) - 1; i >= 0; i-- {
		pl := portLocations[i]
		cmd := fmt.Sprintf("sudo iptables -t nat -D %s %d", pl.Chain, pl.Index)
		task.ExecuteStringTask(cmd)
	}
	return nil
}

func RunLocal(port int) error {
	natOut, err := getNatRules()
	if err != nil {
		return err
	}

	// Exit if we already have iptable rules setup for this port
	portLocations := parseNatRules(natOut, port)
	if len(portLocations) > 0 {
		return nil
	}

	// Masquerade is needed but only once. Does not seem to cause issues if it is not removed.
	if !parseIsMasquerade(natOut) {
		_, err := task.ExecuteStringTask("sudo iptables -t nat -A POSTROUTING -j MASQUERADE")
		if err != nil {
			return err
		}
	}

	for _, ip := range proxy_ips {
		prerouting := fmt.Sprintf("sudo iptables -t nat -A PREROUTING --dst %s -p tcp --dport %d -j DNAT --to-destination %s:%d", ip, port, host_ip, port)
		output := fmt.Sprintf("sudo iptables -t nat -A OUTPUT --dst %s -p tcp --dport %d -j DNAT --to-destination %s:%d", ip, port, host_ip, port)
		task.ExecuteStringTask(prerouting)
		task.ExecuteStringTask(output)
	}
	return nil
}
