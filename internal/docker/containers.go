package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/presselam/yadc/internal/bubble"
	"log"
	"reflect"
	"strconv"
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

func ContainerInspect(id string) (Results, error) {
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

	inspect, err := docker.ContainerInspect(context.Background(), id)
	if err != nil {
		return retval, err
	}

	rows := printObject(inspect.ContainerJSONBase, 0)
	for _, row := range rows {
		retval.Data = append(retval.Data, row)
		for i, val := range row {
			if len(val) > retval.Width[i] {
				retval.Width[i] = len(val)
			}
		}
	}

	//	log.Printf("inspect.Mounts:[%v]", inspect.Mounts)
	//	log.Printf("inspect.Config:[%v]", inspect.Config)
	//	log.Printf("inspect.NetworkSettings:[%v]", inspect.NetworkSettings)
	//	log.Printf("inspect.ImageManifestDescriptor:[%v]", inspect.ImageManifestDescriptor)

	return retval, nil
}

func printObject(obj any, depth int) []bubble.Row {
	var retval []bubble.Row

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := t.Field(i).Name
		//    log.Printf("fld:[%d][%d] => %s:%v", depth, i, name, field.Interface())

		value := convertValue(name, field, depth)
		for _, ln := range value {
			retval = append(retval, ln)
		}
	}
	return retval
}

func convertValue(name string, field reflect.Value, depth int) []bubble.Row {
	padd := fmt.Sprintf("%"+strconv.Itoa(4*depth)+"s", "")
	var retval []bubble.Row
	switch field.Kind() {
	case reflect.String:
		row := bubble.Row{padd + name, field.String()}
		retval = append(retval, row)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		row := bubble.Row{padd + name, strconv.FormatInt(field.Int(), 10)}
		retval = append(retval, row)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		row := bubble.Row{padd + name, strconv.FormatUint(field.Uint(), 10)}
		retval = append(retval, row)
	case reflect.Slice, reflect.Array:
		var sliceValues []string
		value := ""
		if field.Len() > 0 {
			for j := 0; j < field.Len(); j++ {
				sliceValues = append(sliceValues, fmt.Sprint(field.Index(j).Interface()))
			}
			value = "[" + fmt.Sprint(strings.Join(sliceValues, ", ")) + "]"
		}
		row := bubble.Row{padd + name, value}
		retval = append(retval, row)
	case reflect.Pointer:
		if !field.IsNil() {
			retval = append(retval, bubble.Row{padd + name, ""})
			lines := printObject(field.Interface(), depth+1)
			for _, ln := range lines {
				retval = append(retval, ln)
			}
		}
	case reflect.Struct:
		retval = append(retval, bubble.Row{padd + name, ""})
		lines := printObject(field.Interface(), depth+1)
		for _, ln := range lines {
			retval = append(retval, ln)
		}
	case reflect.Bool:
		value := fmt.Sprintf("%v", field.Bool())
		row := bubble.Row{padd + name, value}
		retval = append(retval, row)
	case reflect.Map:
		retval = append(retval, bubble.Row{padd + name + "{", ""})
		keys := field.MapKeys()
		for _, key := range keys {
			log.Printf("map[%v] => [%v]", key, field.MapIndex(key))
			retval = append(retval, bubble.Row{padd + key.String(), ""})
			lines := convertValue(name, field.MapIndex(key), depth+1)
			for _, ln := range lines {
				retval = append(retval, ln)
			}
		}
		retval = append(retval, bubble.Row{padd + "}", ""})
	default:
		value := fmt.Sprintf("TODO: [%v]", field.Kind())
		row := bubble.Row{padd + name, value}
		retval = append(retval, row)
	}
	return retval
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
