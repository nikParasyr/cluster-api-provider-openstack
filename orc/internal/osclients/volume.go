/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package osclients

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type VolumeClient interface {
	ListVolumes(opts volumes.ListOptsBuilder) ([]volumes.Volume, error)
	CreateVolume(opts volumes.CreateOptsBuilder) (*volumes.Volume, error)
	DeleteVolume(volumeID string, opts volumes.DeleteOptsBuilder) error
	GetVolume(volumeID string) (*volumes.Volume, error)
}

type volumeClient struct{ client *gophercloud.ServiceClient }

// NewVolumeClient returns a new cinder client.
func NewVolumeClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (VolumeClient, error) {
	volume, err := openstack.NewBlockStorageV3(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create volume service client: %v", err)
	}

	return &volumeClient{volume}, nil
}

func (c volumeClient) ListVolumes(opts volumes.ListOptsBuilder) ([]volumes.Volume, error) {
	pages, err := volumes.List(c.client, opts).AllPages(context.TODO())
	if err != nil {
		return nil, err
	}
	return volumes.ExtractVolumes(pages)
}

func (c volumeClient) CreateVolume(opts volumes.CreateOptsBuilder) (*volumes.Volume, error) {
	return volumes.Create(context.TODO(), c.client, opts, nil).Extract()
}

func (c volumeClient) DeleteVolume(volumeID string, opts volumes.DeleteOptsBuilder) error {
	return volumes.Delete(context.TODO(), c.client, volumeID, opts).ExtractErr()
}

func (c volumeClient) GetVolume(volumeID string) (*volumes.Volume, error) {
	return volumes.Get(context.TODO(), c.client, volumeID).Extract()
}

type volumeErrorClient struct{ error }

// NewVolumeErrorClient returns a VolumeClient in which every method returns the given error.
func NewVolumeErrorClient(e error) VolumeClient {
	return volumeErrorClient{e}
}

func (e volumeErrorClient) ListVolumes(_ volumes.ListOptsBuilder) ([]volumes.Volume, error) {
	return nil, e.error
}

func (e volumeErrorClient) CreateVolume(_ volumes.CreateOptsBuilder) (*volumes.Volume, error) {
	return nil, e.error
}

func (e volumeErrorClient) DeleteVolume(_ string, _ volumes.DeleteOptsBuilder) error {
	return e.error
}

func (e volumeErrorClient) GetVolume(_ string) (*volumes.Volume, error) {
	return nil, e.error
}
