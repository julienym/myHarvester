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
	"context"
	"time"

	v1beta1 "github.com/kubernetes-csi/external-snapshotter/v2/pkg/apis/volumesnapshot/v1beta1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type VolumeSnapshotClassHandler func(string, *v1beta1.VolumeSnapshotClass) (*v1beta1.VolumeSnapshotClass, error)

type VolumeSnapshotClassController interface {
	generic.ControllerMeta
	VolumeSnapshotClassClient

	OnChange(ctx context.Context, name string, sync VolumeSnapshotClassHandler)
	OnRemove(ctx context.Context, name string, sync VolumeSnapshotClassHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() VolumeSnapshotClassCache
}

type VolumeSnapshotClassClient interface {
	Create(*v1beta1.VolumeSnapshotClass) (*v1beta1.VolumeSnapshotClass, error)
	Update(*v1beta1.VolumeSnapshotClass) (*v1beta1.VolumeSnapshotClass, error)

	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1beta1.VolumeSnapshotClass, error)
	List(opts metav1.ListOptions) (*v1beta1.VolumeSnapshotClassList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.VolumeSnapshotClass, err error)
}

type VolumeSnapshotClassCache interface {
	Get(name string) (*v1beta1.VolumeSnapshotClass, error)
	List(selector labels.Selector) ([]*v1beta1.VolumeSnapshotClass, error)

	AddIndexer(indexName string, indexer VolumeSnapshotClassIndexer)
	GetByIndex(indexName, key string) ([]*v1beta1.VolumeSnapshotClass, error)
}

type VolumeSnapshotClassIndexer func(obj *v1beta1.VolumeSnapshotClass) ([]string, error)

type volumeSnapshotClassController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewVolumeSnapshotClassController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) VolumeSnapshotClassController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &volumeSnapshotClassController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromVolumeSnapshotClassHandlerToHandler(sync VolumeSnapshotClassHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1beta1.VolumeSnapshotClass
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1beta1.VolumeSnapshotClass))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *volumeSnapshotClassController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1beta1.VolumeSnapshotClass))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateVolumeSnapshotClassDeepCopyOnChange(client VolumeSnapshotClassClient, obj *v1beta1.VolumeSnapshotClass, handler func(obj *v1beta1.VolumeSnapshotClass) (*v1beta1.VolumeSnapshotClass, error)) (*v1beta1.VolumeSnapshotClass, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *volumeSnapshotClassController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *volumeSnapshotClassController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *volumeSnapshotClassController) OnChange(ctx context.Context, name string, sync VolumeSnapshotClassHandler) {
	c.AddGenericHandler(ctx, name, FromVolumeSnapshotClassHandlerToHandler(sync))
}

func (c *volumeSnapshotClassController) OnRemove(ctx context.Context, name string, sync VolumeSnapshotClassHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromVolumeSnapshotClassHandlerToHandler(sync)))
}

func (c *volumeSnapshotClassController) Enqueue(name string) {
	c.controller.Enqueue("", name)
}

func (c *volumeSnapshotClassController) EnqueueAfter(name string, duration time.Duration) {
	c.controller.EnqueueAfter("", name, duration)
}

func (c *volumeSnapshotClassController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *volumeSnapshotClassController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *volumeSnapshotClassController) Cache() VolumeSnapshotClassCache {
	return &volumeSnapshotClassCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *volumeSnapshotClassController) Create(obj *v1beta1.VolumeSnapshotClass) (*v1beta1.VolumeSnapshotClass, error) {
	result := &v1beta1.VolumeSnapshotClass{}
	return result, c.client.Create(context.TODO(), "", obj, result, metav1.CreateOptions{})
}

func (c *volumeSnapshotClassController) Update(obj *v1beta1.VolumeSnapshotClass) (*v1beta1.VolumeSnapshotClass, error) {
	result := &v1beta1.VolumeSnapshotClass{}
	return result, c.client.Update(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *volumeSnapshotClassController) Delete(name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), "", name, *options)
}

func (c *volumeSnapshotClassController) Get(name string, options metav1.GetOptions) (*v1beta1.VolumeSnapshotClass, error) {
	result := &v1beta1.VolumeSnapshotClass{}
	return result, c.client.Get(context.TODO(), "", name, result, options)
}

func (c *volumeSnapshotClassController) List(opts metav1.ListOptions) (*v1beta1.VolumeSnapshotClassList, error) {
	result := &v1beta1.VolumeSnapshotClassList{}
	return result, c.client.List(context.TODO(), "", result, opts)
}

func (c *volumeSnapshotClassController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), "", opts)
}

func (c *volumeSnapshotClassController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*v1beta1.VolumeSnapshotClass, error) {
	result := &v1beta1.VolumeSnapshotClass{}
	return result, c.client.Patch(context.TODO(), "", name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type volumeSnapshotClassCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *volumeSnapshotClassCache) Get(name string) (*v1beta1.VolumeSnapshotClass, error) {
	obj, exists, err := c.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1beta1.VolumeSnapshotClass), nil
}

func (c *volumeSnapshotClassCache) List(selector labels.Selector) (ret []*v1beta1.VolumeSnapshotClass, err error) {

	err = cache.ListAll(c.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.VolumeSnapshotClass))
	})

	return ret, err
}

func (c *volumeSnapshotClassCache) AddIndexer(indexName string, indexer VolumeSnapshotClassIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1beta1.VolumeSnapshotClass))
		},
	}))
}

func (c *volumeSnapshotClassCache) GetByIndex(indexName, key string) (result []*v1beta1.VolumeSnapshotClass, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1beta1.VolumeSnapshotClass, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1beta1.VolumeSnapshotClass))
	}
	return result, nil
}