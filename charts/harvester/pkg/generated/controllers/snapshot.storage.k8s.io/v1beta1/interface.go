/*
Copyright 2022 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/kubernetes-csi/external-snapshotter/v2/pkg/apis/volumesnapshot/v1beta1"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/schemes"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	schemes.Register(v1beta1.AddToScheme)
}

type Interface interface {
	VolumeSnapshot() VolumeSnapshotController
	VolumeSnapshotClass() VolumeSnapshotClassController
	VolumeSnapshotContent() VolumeSnapshotContentController
}

func New(controllerFactory controller.SharedControllerFactory) Interface {
	return &version{
		controllerFactory: controllerFactory,
	}
}

type version struct {
	controllerFactory controller.SharedControllerFactory
}

func (c *version) VolumeSnapshot() VolumeSnapshotController {
	return NewVolumeSnapshotController(schema.GroupVersionKind{Group: "snapshot.storage.k8s.io", Version: "v1beta1", Kind: "VolumeSnapshot"}, "volumesnapshots", true, c.controllerFactory)
}
func (c *version) VolumeSnapshotClass() VolumeSnapshotClassController {
	return NewVolumeSnapshotClassController(schema.GroupVersionKind{Group: "snapshot.storage.k8s.io", Version: "v1beta1", Kind: "VolumeSnapshotClass"}, "volumesnapshotclasses", false, c.controllerFactory)
}
func (c *version) VolumeSnapshotContent() VolumeSnapshotContentController {
	return NewVolumeSnapshotContentController(schema.GroupVersionKind{Group: "snapshot.storage.k8s.io", Version: "v1beta1", Kind: "VolumeSnapshotContent"}, "volumesnapshotcontents", false, c.controllerFactory)
}