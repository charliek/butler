package service

import (
	"fmt"
	"github.com/charliek/butler/iptables"
	"github.com/charliek/butler/task"
	log "github.com/ngmoco/timber"
)

type ButlerService struct {
	Name        string
	Display     string
	ServiceType string
	Port        int
}

func (s *ButlerService) Stop() error {
	cmd := fmt.Sprintf("sudo service %s stop", s.Name)
	stdout, err := task.ExecuteStringTask(cmd)
	log.Info("Service output: %s", stdout)
	return err
}

func (s *ButlerService) Start() error {
	cmd := fmt.Sprintf("sudo service %s start", s.Name)
	stdout, err := task.ExecuteStringTask(cmd)
	log.Info("Service output: %s", stdout)
	return err
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
