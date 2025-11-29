package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/presselam/yadc/internal/logger"
	"strconv"
)

func Volumes() (Results, error) {
	logger.Trace()
	retval := Results{
		[]string{"ID", "Driver", "Mount", "Created", "Scope", "Size"},
		[][]string{},
		[]int{0, 0, 0, 0, 0, 0},
	}

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return retval, err
	}
	defer docker.Close()

	resp, err := docker.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		return retval, err
	}

	for _, vol := range resp.Volumes {
		var size string
		if vol.UsageData != nil {
			size = strconv.FormatInt(vol.UsageData.Size, 10)
		}
		logger.Debug(vol.UsageData)
		row := []string{
			vol.Name,
			vol.Driver,
			vol.Mountpoint,
			vol.CreatedAt,
			vol.Scope,
			size,
		}
		retval.Data = append(retval.Data, row)

		for i, val := range row {
			retval.Width[i] = min(max(len(val), retval.Width[i]), 20)
		}
	}

	return retval, nil
}

func VolumeDelete(id string) (string, error) {
	logger.Trace(id)

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	defer docker.Close()

	err = docker.VolumeRemove(context.Background(), id, false)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Removed Volume: %s", id), nil
}

func VolumesPrune(id string) (string, error) {
	logger.Trace(id)

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	defer docker.Close()

	filters := filters.NewArgs()
	report, err := docker.VolumesPrune(context.Background(), filters)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	retval := fmt.Sprintf("Removed: %d Images\nTotal reclaimed space: %d", len(report.VolumesDeleted), report.SpaceReclaimed)
	logger.Debug(retval)

	return retval, nil
}

func VolumeInspect(id string) (Results, error) {
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

	inspect, err := docker.VolumeInspect(context.Background(), id)
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

	return retval, nil
}
