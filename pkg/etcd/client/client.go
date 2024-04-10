package client

import (
	"etcd-operator/pkg/etcd"
)

// Member contains member info including v2 and v3
type Member struct {
	ID         string
	Name       string
	PeerURLs   []string
	ClientURLs []string
	Version    string
	IsLearner  bool
	Leader     string
}

type VersionClient interface {
	MemberList() ([]Member, error)
	Status(endpoint string) (*Member, error)
	Close()
}

type VersionContext struct {
	Config *etcd.ClientConfig
}
