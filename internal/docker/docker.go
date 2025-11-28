package docker

import (
	"context"
	"github.com/docker/docker/client"
)

type ServerInfo struct {
	ID            string
	Running       int
	Paused        int
	Stopped       int
	Images        int
	Name          string
	ServerVersion string
	ClientVersion string
}

type Results struct {
	Columns []string
	Data    [][]string
	Width   []int
}

func Info() (ServerInfo, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer docker.Close()

	info, err := docker.Info(context.Background())
	retval := ServerInfo{
		info.ID,
		info.ContainersRunning,
		info.ContainersPaused,
		info.ContainersStopped,
		info.Images,
		info.Name,
		info.ServerVersion,
		docker.ClientVersion(),
	}

	return retval, err
}
