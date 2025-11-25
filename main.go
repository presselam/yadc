package main

import (
	"flag"
	"github.com/presselam/yadc/internal/monitor"
)

func main() {
	container := flag.Bool("containers", false, "start the monitor in container mode")
	image := flag.Bool("images", false, "start the monitor in images mode")
	volume := flag.Bool("volumes", false, "start the monitor in volume mode")
	flag.Parse()

	var mode string
	switch {
	case *container:
		mode = monitor.ContainerMode
	case *image:
		mode = monitor.ImageMode
	case *volume:
		mode = monitor.VolumeMode
	default:
		mode = monitor.ContainerMode
	}

	monitor.Show(mode)
}
