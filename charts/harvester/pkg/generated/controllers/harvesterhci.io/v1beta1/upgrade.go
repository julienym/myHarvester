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

	v1beta1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
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

type UpgradeHandler func(string, *v1beta1.Upgrade) (*v1beta1.Upgrade, error)

type UpgradeController interface {
	generic.ControllerMeta
	UpgradeClient

	OnChange(ctx context.Context, name string, sync UpgradeHandler)
	OnRemove(ctx context.Context, name string, sync UpgradeHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() UpgradeCache
}

type UpgradeClient interface {
	Create(*v1beta1.Upgrade) (*v1beta1.Upgrade, error)
	Update(*v1beta1.Upgrade) (*v1beta1.Upgrade, error)
	UpdateStatus(*v1beta1.Upgrade) (*v1beta1.Upgrade, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1beta1.Upgrade, error)
	List(namespace string, opts metav1.ListOptions) (*v1beta1.UpgradeList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.Upgrade, err error)
}

type UpgradeCache interface {
	Get(namespace, name string) (*v1beta1.Upgrade, error)
	List(namespace string, selector labels.Selector) ([]*v1beta1.Upgrade, error)

	AddIndexer(indexName string, indexer UpgradeIndexer)
	GetByIndex(indexName, key string) ([]*v1beta1.Upgrade, error)
}

type UpgradeIndexer func(obj *v1beta1.Upgrade) ([]string, error)

type upgradeController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewUpgradeController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) UpgradeController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &upgradeController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromUpgradeHandlerToHandler(sync UpgradeHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1beta1.Upgrade
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1beta1.Upgrade))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *upgradeController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1beta1.Upgrade))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateUpgradeDeepCopyOnChange(client UpgradeClient, obj *v1beta1.Upgrade, handler func(obj *v1beta1.Upgrade) (*v1beta1.Upgrade, error)) (*v1beta1.Upgrade, error) {
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

func (c *upgradeController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *upgradeController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *upgradeController) OnChange(ctx context.Context, name string, sync UpgradeHandler) {
	c.AddGenericHandler(ctx, name, FromUpgradeHandlerToHandler(sync))
}

func (c *upgradeController) OnRemove(ctx context.Context, name string, sync UpgradeHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromUpgradeHandlerToHandler(sync)))
}

func (c *upgradeController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *upgradeController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *upgradeController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *upgradeController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *upgradeController) Cache() UpgradeCache {
	return &upgradeCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *upgradeController) Create(obj *v1beta1.Upgrade) (*v1beta1.Upgrade, error) {
	result := &v1beta1.Upgrade{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *upgradeController) Update(obj *v1beta1.Upgrade) (*v1beta1.Upgrade, error) {
	result := &v1beta1.Upgrade{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *upgradeController) UpdateStatus(obj *v1beta1.Upgrade) (*v1beta1.Upgrade, error) {
	result := &v1beta1.Upgrade{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *upgradeController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *upgradeController) Get(namespace, name string, options metav1.GetOptions) (*v1beta1.Upgrade, error) {
	result := &v1beta1.Upgrade{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *upgradeController) List(namespace string, opts metav1.ListOptions) (*v1beta1.UpgradeList, error) {
	result := &v1beta1.UpgradeList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *upgradeController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *upgradeController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1beta1.Upgrade, error) {
	result := &v1beta1.Upgrade{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type upgradeCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *upgradeCache) Get(namespace, name string) (*v1beta1.Upgrade, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1beta1.Upgrade), nil
}

func (c *upgradeCache) List(namespace string, selector labels.Selector) (ret []*v1beta1.Upgrade, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.Upgrade))
	})

	return ret, err
}

func (c *upgradeCache) AddIndexer(indexName string, indexer UpgradeIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1beta1.Upgrade))
		},
	}))
}

func (c *upgradeCache) GetByIndex(indexName, key string) (result []*v1beta1.Upgrade, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1beta1.Upgrade, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1beta1.Upgrade))
	}
	return result, nil
}

type UpgradeStatusHandler func(obj *v1beta1.Upgrade, status v1beta1.UpgradeStatus) (v1beta1.UpgradeStatus, error)

type UpgradeGeneratingHandler func(obj *v1beta1.Upgrade, status v1beta1.UpgradeStatus) ([]runtime.Object, v1beta1.UpgradeStatus, error)

func RegisterUpgradeStatusHandler(ctx context.Context, controller UpgradeController, condition condition.Cond, name string, handler UpgradeStatusHandler) {
	statusHandler := &upgradeStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromUpgradeHandlerToHandler(statusHandler.sync))
}

func RegisterUpgradeGeneratingHandler(ctx context.Context, controller UpgradeController, apply apply.Apply,
	condition condition.Cond, name string, handler UpgradeGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &upgradeGeneratingHandler{
		UpgradeGeneratingHandler: handler,
		apply:                    apply,
		name:                     name,
		gvk:                      controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterUpgradeStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type upgradeStatusHandler struct {
	client    UpgradeClient
	condition condition.Cond
	handler   UpgradeStatusHandler
}

func (a *upgradeStatusHandler) sync(key string, obj *v1beta1.Upgrade) (*v1beta1.Upgrade, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type upgradeGeneratingHandler struct {
	UpgradeGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *upgradeGeneratingHandler) Remove(key string, obj *v1beta1.Upgrade) (*v1beta1.Upgrade, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1beta1.Upgrade{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *upgradeGeneratingHandler) Handle(obj *v1beta1.Upgrade, status v1beta1.UpgradeStatus) (v1beta1.UpgradeStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.UpgradeGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
