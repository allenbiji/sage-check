package checks

import (
	"fmt"
	"net"

	"github.com/allenbiji/clone-sage/internal/model"
	"github.com/allenbiji/clone-sage/internal/registry"
)

type PortFreeCheck struct {
	Port string
}

//execute method for the port free check
func (p *PortFreeCheck) Execute() error {
	listen, err := net.Listen("tcp", ":" + p.Port)
	if err != nil {
		return fmt.Errorf("The port '%s' is  not free", p.Port)
	}

	defer listen.Close()

	return nil
}

//build a factory for the port free check
func buildPortFreeCheck(cfg model.CheckConfig) (registry.Check, error) {
	port, ok := cfg.Options["port"]
	if !ok || port == "" {
		return nil, fmt.Errorf("The port_free check requires a 'port' option")
	}

	return &PortFreeCheck{
		Port: port,
	}, nil
}

//registers the check in the registry
func init() {
	registry.Register(model.TypePortFree, buildPortFreeCheck)
}