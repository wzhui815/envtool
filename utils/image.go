package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/opencontainers/go-digest"

	"github.com/docker/docker/image"
	v1 "github.com/docker/docker/image/v1"
	"github.com/docker/docker/layer"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ListContainers list the containers local host
func ListContainers() {
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
	for _, imageSummary := range images {
		fmt.Printf("%s %s\n", imageSummary.ID, imageSummary.RepoTags[0])

		imageSpec, _, err := cli.ImageInspectWithRaw(context.Background(), imageSummary.ID)
		if err != nil {
			panic(err)
		}

		var chainID layer.ChainID
		var parent digest.Digest
		for i, layerID := range imageSpec.RootFS.Layers {
			if chainID == "" {
				chainID = layer.ChainID(layerID)
			} else {
				chainID = layer.ChainID(digest.FromString(chainID.String() + " " + layerID))
			}

			img := image.V1Image{
				Created: time.Unix(0, 0),
			}
			if i == len(imageSpec.RootFS.Layers)-1 {
				fmt.Println("=====================================================", i)
				img = sepc2V1Image(imageSpec)
			}
			// var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
			// img.Created, _ = time.ParseInLocation("1/2/2006 15:04:05", "1/1/1970 00:00:00", cstSh)

			//fmt.Println("SH : ", time.Now().In(cstSh).Format("2006-01-02 15:04:05"))

			v1Id, _ := v1.CreateID(img, chainID, parent)
			// v1Id, _ := CreateID(V1Image{}, chainID, parent)
			// v1Id := CreateID(chainID, parent)
			fmt.Printf(">>> ChainID: %s -> saveID: %s\n", chainID, v1Id)

			parent = v1Id
		}

	}

}

// 	types.ImageInspect
// ID              string `json:"Id"`
// RepoTags        []string
// RepoDigests     []string
// Parent          string
// Comment         string
// Created         string
// Container       string
// ContainerConfig *container.Config
// DockerVersion   string
// Author          string
// Config          *container.Config
// Architecture    string
// Variant         string `json:",omitempty"`
// Os              string
// OsVersion       string `json:",omitempty"`
// Size            int64
// VirtualSize     int64
// GraphDriver     GraphDriverData
// RootFS          RootFS
// Metadata        ImageMetadata

func sepc2V1Image(spec types.ImageInspect) image.V1Image {
	tm, _ := time.Parse(time.RFC3339, spec.Created)

	v1Img := image.V1Image{
		ID:              spec.ID,
		Parent:          spec.Parent,
		Comment:         spec.Comment,
		Created:         tm,
		Container:       spec.Container,
		ContainerConfig: *spec.ContainerConfig,
		DockerVersion:   spec.DockerVersion,
		Author:          spec.Author,
		Config:          spec.Config,
		Architecture:    spec.Architecture,
		Variant:         spec.Variant,
		OS:              spec.Os,
		// Size:            spec.Size,
	}
	return v1Img
}
