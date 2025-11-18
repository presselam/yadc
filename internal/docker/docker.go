package docker

import (
  "context"
  "github.com/docker/docker/api/types/container"
//  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
)

func Containers() ([]container.Summary, error) {
  docker, err := client.NewClientWithOpts(client.FromEnv)
  if err != nil {
    panic(err)
  }
  defer docker.Close()

  containers, err := docker.ContainerList(context.Background(), container.ListOptions{All: true})
  return containers, err
}

