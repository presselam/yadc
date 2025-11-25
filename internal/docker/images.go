package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"log"
	"slices"
	"strconv"
	"strings"
)

const (
	shaPrefix    = "sha256:"
	imageNone    = "<none>"
	imageMissing = "<missing>"
)

func Images() (Results, error) {
	log.Println("docker.images.Images")
	retval := Results{
		[]string{"ID", "Name", "Contianers", "Size"},
		[][]string{},
		[]int{0, 0, 0, 0},
	}

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return retval, err
	}
	defer docker.Close()

	images, err := docker.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		return retval, err
	}

	for _, img := range images {
		if strings.HasPrefix(img.ID, shaPrefix) {
			img.ID = strings.TrimPrefix(img.ID, shaPrefix)
		}
		img.ID = img.ID[0:8]

		names := img.RepoTags
		if len(names) == 0 {
			names = append(names, imageNone)
		}

		for _, name := range names {
			row := []string{
				img.ID,
				name,
				strconv.FormatInt(img.Containers, 10),
				strconv.FormatInt(img.Size, 10),
			}
			retval.Data = append(retval.Data, row)

			for i, val := range row {
				if len(val) > retval.Width[i] {
					retval.Width[i] = len(val)
				}
			}
		}
	}

	return retval, nil
}

func ImageDelete(id string) (string, error) {
	log.Printf("docker.images.ImageDelete.%s", id)

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	defer docker.Close()

	response, err := docker.ImageRemove(context.Background(), id, image.RemoveOptions{})
	if err != nil {
		return "", err
	}

	var retval string
	for _, img := range response {
		retval = fmt.Sprintf("%s\n%s", img.Untagged, img.Deleted)
	}

	return retval, nil
}

func ImageHistory(id string) (Results, error) {
	log.Printf("docker.images.ImageHistory.%s", id)
	retval := Results{
		[]string{"ID", "Created", "Created By", "Size", "Comment"},
		[][]string{},
		[]int{0, 0, 0, 0, 0},
	}

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return retval, err
	}
	defer docker.Close()

	response, err := docker.ImageHistory(context.Background(), id)
	if err != nil {
		return retval, err
	}

	slices.Reverse(response)

	for _, layer := range response {
		if strings.HasPrefix(layer.ID, shaPrefix) {
			layer.ID = strings.TrimPrefix(layer.ID, shaPrefix)
			layer.ID = layer.ID[0:8]
		} else if layer.ID == imageMissing {
			layer.ID = ""
		}

		row := []string{
			layer.ID,
			strconv.FormatInt(layer.Created, 10),
			layer.CreatedBy,
			strconv.FormatInt(layer.Size, 10),
			layer.Comment,
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

func ImagesPrune(id string) (string, error) {
	log.Println("docker.images.ImagePrune")

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	defer docker.Close()

	report, err := docker.ImagesPrune(context.Background(), filters.Args{})
	if err != nil {
		return "", err
	}

	retval := fmt.Sprintf("Removed: %d Images\nTotal reclaimed space: %d", len(report.ImagesDeleted), report.SpaceReclaimed)

	return retval, nil
}

func ImageInspect(id string) (Results, error) {
	retval := Results{
		[]string{"Name", "Value"},
		[][]string{},
		[]int{0, 0},
	}
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return retval, err
	}
	defer docker.Close()

	inspect, err := docker.ImageInspect(context.Background(), id)
	if err != nil {
		return retval, err
	}

	rows := printObject(inspect, 0)
	for _, row := range rows {
		retval.Data = append(retval.Data, row)
		for i, val := range row {
			if len(val) > retval.Width[i] {
				retval.Width[i] = len(val)
			}
		}
	}

	//  log.Printf("inspect.Mounts:[%v]", inspect.Mounts)
	//  log.Printf("inspect.Config:[%v]", inspect.Config)
	//  log.Printf("inspect.NetworkSettings:[%v]", inspect.NetworkSettings)
	//  log.Printf("inspect.ImageManifestDescriptor:[%v]", inspect.ImageManifestDescriptor)

	return retval, nil
}
