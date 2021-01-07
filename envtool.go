package main

import (
	"context"
	"fmt"

	"github.com/opencontainers/go-digest"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	for _, image := range images {
		fmt.Printf("%s %s\n", image.ID, image.RepoTags[0])

		imageSpec, _, err := cli.ImageInspectWithRaw(context.Background(), image.ID)
		if err != nil {
			panic(err)
		}

		chainID := ""
		for _, layer := range imageSpec.RootFS.Layers {
			if chainID == "" {
				chainID = layer
			} else {
				chainID = digest.FromString(chainID + " " + layer).String()
			}
			fmt.Printf(">>> ChainID: %s\n", chainID)
		}
	}

}
