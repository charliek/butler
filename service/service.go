package service

import (
	"github.com/charliek/butler/iptables"
	log "github.com/ngmoco/timber"
)

type ButlerService struct {
	Name        string
	Display     string
	ServiceType string
	Port        int
}

func (s *ButlerService) RunLocal() error {
	return iptables.RunLocal(s.Port)
}

func (s *ButlerService) RunVagrant() error {
	return iptables.RunVagrant(s.Port)
}

func (s *ButlerService) IsLocal() bool {
	local, err := iptables.IsLocal(s.Port)
	if err != nil {
		log.Warn("Error determining if service is local %v", err)
	}
	return local
}
