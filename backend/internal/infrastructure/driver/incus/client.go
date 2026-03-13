package incus

import "github.com/lxc/incus/v6/client"

func NewClient(cfg *Config) incus.InstanceServer {
	c, err := incus.ConnectIncusUnix(cfg.SocketPath, nil)
	if err != nil {
		panic(err)
	}
	return c
}