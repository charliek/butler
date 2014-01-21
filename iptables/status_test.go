package iptables

import (
	"github.com/bmizerany/assert"
	"testing"
)

var natRules = `Chain PREROUTING (policy ACCEPT)
target     prot opt source               destination         
DNAT       tcp  --  anywhere             192.168.70.4         tcp dpt:8096 to:192.168.70.1:8096
DNAT       tcp  --  anywhere             10.0.2.15            tcp dpt:8096 to:192.168.70.1:8096
DNAT       tcp  --  anywhere             192.168.70.4         tcp dpt:9090 to:192.168.70.1:9090
DNAT       tcp  --  anywhere             10.0.2.15            tcp dpt:9090 to:192.168.70.1:9090
DNAT       tcp  --  anywhere             192.168.70.4         tcp dpt:http-alt to:192.168.70.1:8080
DNAT       tcp  --  anywhere             10.0.2.15            tcp dpt:http-alt to:192.168.70.1:8080

Chain INPUT (policy ACCEPT)
target     prot opt source               destination         

Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination         
DNAT       tcp  --  anywhere             192.168.70.4         tcp dpt:8096 to:192.168.70.1:8096
DNAT       tcp  --  anywhere             10.0.2.15            tcp dpt:8096 to:192.168.70.1:8096
DNAT       tcp  --  anywhere             192.168.70.4         tcp dpt:9090 to:192.168.70.1:9090
DNAT       tcp  --  anywhere             10.0.2.15            tcp dpt:9090 to:192.168.70.1:9090
DNAT       tcp  --  anywhere             192.168.70.4         tcp dpt:http-alt to:192.168.70.1:8080
DNAT       tcp  --  anywhere             10.0.2.15            tcp dpt:http-alt to:192.168.70.1:8080

Chain POSTROUTING (policy ACCEPT)
target     prot opt source               destination         
MASQUERADE  all  --  anywhere             anywhere      
`

func assertPortLocation(t *testing.T, p *PortLocation, chain string, loc int) {
	assert.Equal(t, p.Chain, chain)
	assert.Equal(t, p.Index, loc)
}

func TestParseNatRules(t *testing.T) {
	portLocations := parseNatRules(natRules, 9090)
	assert.Equal(t, 4, len(portLocations))
	assertPortLocation(t, portLocations[0], "PREROUTING", 3)
	assertPortLocation(t, portLocations[1], "PREROUTING", 4)
	assertPortLocation(t, portLocations[2], "OUTPUT", 3)
	assertPortLocation(t, portLocations[3], "OUTPUT", 4)
}

func TestParseNatRulesRenamed(t *testing.T) {
	portLocations := parseNatRules(natRules, 8080)
	assert.Equal(t, 4, len(portLocations))
	assertPortLocation(t, portLocations[0], "PREROUTING", 5)
	assertPortLocation(t, portLocations[1], "PREROUTING", 6)
	assertPortLocation(t, portLocations[2], "OUTPUT", 5)
	assertPortLocation(t, portLocations[3], "OUTPUT", 6)
}

// Commenting out due to sudo need
// func TestGetNatRules(t *testing.T) {
// 	s, err := getNatRules()
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, "", s)
// }
