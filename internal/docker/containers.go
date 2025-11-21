package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"strings"
)

func Containers() (Results, error) {
	retval := Results{
		[]string{"ID", "Name", "Image", "State", "Ports"},
		[][]string{},
		[]int{0, 0, 0, 0, 0},
	}

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return retval, err
	}
	defer docker.Close()

	containers, err := docker.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return retval, err
	}

	for _, cont := range containers {
		var row []string
		if len(cont.Names) == 0 {
			row = []string{
				cont.ID[0:8],
				"<none>",
				cont.Image,
				cont.State,
				displayPorts(cont.Ports),
			}
		} else {
			for _, name := range cont.Names {
				row = []string{
					cont.ID[0:8],
					name[1:],
					cont.Image,
					cont.State,
					displayPorts(cont.Ports),
				}
			}
		}
		retval.Data = append(retval.Data, row)

		for i, val := range row {
			if len(val) > retval.Width[i] {
				retval.Width[i] = len(val)
			}

		}
	}

	return retval, nil
}

func displayPorts(ports []container.Port) string {
	var retval string

	uniq := make(map[string]bool)
	for _, p := range ports {
		lbl := fmt.Sprintf("%d:%d", p.PublicPort, p.PrivatePort)
		uniq[lbl] = true
	}

	keys := []string{}
	for key := range uniq {
		keys = append(keys, key)
	}

	retval += strings.Join(keys, ",")
	return retval
}
