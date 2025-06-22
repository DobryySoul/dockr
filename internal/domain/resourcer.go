package domain

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
)

type UnusedResources struct {
	Images     []image.Summary
	Containers []container.Summary
	Volumes    []volume.Volume
	// Networks   []types.NetworkResource
}

func (r *UnusedResources) TotalSize() float64 {
	var total float64
	for _, img := range r.Images {
		total += float64(img.Size)
	}
	return total
}

func (r *UnusedResources) ImagesSize() float64 {
	var total float64
	for _, img := range r.Images {
		total += float64(img.Size)
	}
	return total
}
