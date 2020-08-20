/*
Copyright 2018 The Kubernetes Authors.

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

package linodebs

import (
	csi "github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/wrappers"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LinodeIdentityServer struct {
	Driver *LinodeDriver
}

// GetPluginInfo(context.Context, *GetPluginInfoRequest) (*GetPluginInfoResponse, error)
func (linodeIdentity *LinodeIdentityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	glog.V(5).Infof("Using default GetPluginInfo")

	if linodeIdentity.Driver.name == "" {
		return nil, status.Error(codes.Unavailable, "Driver name not configured")
	}

	return &csi.GetPluginInfoResponse{
		Name:          linodeIdentity.Driver.name,
		VendorVersion: linodeIdentity.Driver.vendorVersion,
	}, nil
}

func (linodeIdentity *LinodeIdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	glog.V(5).Infof("Using default GetPluginCapabilities")
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS,
					},
				},
			},
			{
				Type: &csi.PluginCapability_VolumeExpansion_{
					VolumeExpansion: &csi.PluginCapability_VolumeExpansion{
						Type: csi.PluginCapability_VolumeExpansion_ONLINE,
					},
				},
			},
		},
	}, nil
}

func (linodeIdentity *LinodeIdentityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	glog.V(4).Infof("Probe called with args: %#v", req)
	linodeIdentity.Driver.readyMu.Lock()
	defer linodeIdentity.Driver.readyMu.Unlock()

	return &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{
			Value: true,
		},
	}, nil
}
