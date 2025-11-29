package main

import (
	"github.com/presselam/yadc/internal/monitor"
	flag "github.com/spf13/pflag"
)

func main() {
	container := flag.BoolP("containers", "c", false, "start the monitor in container mode")
	image := flag.BoolP("images", "i", false, "start the monitor in images mode")
	volume := flag.BoolP("volumes", "v", false, "start the monitor in volume mode")
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
