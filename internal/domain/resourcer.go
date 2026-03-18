package domain

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
)

type UnusedResources struct {
	Images     []*image.Summary
	Containers []*container.Summary
	Volumes    []*volume.Volume
	Networks   []*network.Summary
}

func (ur *UnusedResources) ContainersSize() float64 {
	var total float64
	for _, cont := range ur.Containers {
		total += float64(cont.SizeRw)
	}
	return total
}

func (ur *UnusedResources) TotalSize() float64 {
	var total float64
	for _, img := range ur.Images {
		total += float64(img.Size)
	}
	for _, vol := range ur.Volumes {
		if vol.UsageData != nil {
			total += float64(vol.UsageData.Size)
		}
	}
	for _, cont := range ur.Containers {
		total += float64(cont.SizeRw)
	}
	return total
}

func (ur *UnusedResources) ImagesSize() float64 {
	var total float64
	for _, img := range ur.Images {
		total += float64(img.Size)
	}
	return total
}

func (ur *UnusedResources) VolumesSize() float64 {
	var total float64
	for _, vol := range ur.Volumes {
		if vol.UsageData != nil {
			total += float64(vol.UsageData.Size)
		}
	}
	return total
}

func (ur *UnusedResources) TotalCount() int {
	return len(ur.Images) + len(ur.Containers) + len(ur.Volumes) + len(ur.Networks)
}

func (ur *UnusedResources) IsEmpty() bool {
	return len(ur.Images) == 0 && len(ur.Containers) == 0 && len(ur.Volumes) == 0 && len(ur.Networks) == 0
}
