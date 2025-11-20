package docker

import (
"fmt"
 "strconv"
  "context"
  "github.com/docker/docker/api/types/image"
  "github.com/docker/docker/client"
)

func Images() (Results, error) {
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
    var row []string
    if len(img.RepoTags) == 0 {
        row = []string{
          img.ID[6:14],
          "<none>",
          strconv.FormatInt(img.Containers, 10),
          strconv.FormatInt(img.Size, 10),
        }
    }else{
      for _, tag := range img.RepoTags{
        row = []string{
          img.ID[7:15],
          tag,
          strconv.FormatInt(img.Containers, 10),
          strconv.FormatInt(img.Size, 10),
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

func ImageDelete(id string) (string, error) {
  docker, err := client.NewClientWithOpts(client.FromEnv)
  if err != nil {
    return "", err
  }
  defer docker.Close()

  response, err :=  docker.ImageRemove(context.Background(), id, image.RemoveOptions{})
  if err != nil {
    return "", err
  }

  var retval string;
  for _, img := range response {
  retval = fmt.Sprintf("%s\n%s", img.Untagged, img.Deleted)
  }

  return retval, nil
}
