package deploy

import (
	"github.com/flant/logboek"
	"github.com/flant/werf/pkg/docker_registry"
)

type ImageInfoGetterStub struct {
	Name              string
	Tag               string
	ImagesRepoManager ImagesRepoManager
}

func (d *ImageInfoGetterStub) IsNameless() bool {
	return d.Name == ""
}

func (d *ImageInfoGetterStub) GetName() string {
	return d.Name
}

func (d *ImageInfoGetterStub) GetImageName() string {
	return d.ImagesRepoManager.ImageRepoWithTag(d.Name, d.Tag)
}

func (d *ImageInfoGetterStub) GetImageId() (string, error) {
	return docker_registry.ImageId(d.GetImageName())
}

type ImageInfo struct {
	Name              string
	WithoutRegistry   bool
	ImagesRepoManager ImagesRepoManager
	Tag               string
}

func (d *ImageInfo) IsNameless() bool {
	return d.Name == ""
}

func (d *ImageInfo) GetName() string {
	return d.Name
}

func (d *ImageInfo) GetImageName() string {
	return d.ImagesRepoManager.ImageRepoWithTag(d.Name, d.Tag)
}

func (d *ImageInfo) GetImageId() (string, error) {
	if d.WithoutRegistry {
		return "", nil
	}

	imageName := d.GetImageName()

	res, err := docker_registry.ImageId(imageName)
	if err != nil {
		logboek.LogErrorF("WARNING: Getting image %s id failed: %s\n", imageName, err)
		return "", nil
	}

	return res, nil
}
