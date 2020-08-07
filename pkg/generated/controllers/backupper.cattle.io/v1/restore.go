/*
Copyright 2020 Rancher Labs, Inc.

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

package v1

import (
	"context"
	"time"

	v1 "github.com/mrajashree/backup/pkg/apis/backupper.cattle.io/v1"
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

type RestoreHandler func(string, *v1.Restore) (*v1.Restore, error)

type RestoreController interface {
	generic.ControllerMeta
	RestoreClient

	OnChange(ctx context.Context, name string, sync RestoreHandler)
	OnRemove(ctx context.Context, name string, sync RestoreHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() RestoreCache
}

type RestoreClient interface {
	Create(*v1.Restore) (*v1.Restore, error)
	Update(*v1.Restore) (*v1.Restore, error)
	UpdateStatus(*v1.Restore) (*v1.Restore, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Restore, error)
	List(namespace string, opts metav1.ListOptions) (*v1.RestoreList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Restore, err error)
}

type RestoreCache interface {
	Get(namespace, name string) (*v1.Restore, error)
	List(namespace string, selector labels.Selector) ([]*v1.Restore, error)

	AddIndexer(indexName string, indexer RestoreIndexer)
	GetByIndex(indexName, key string) ([]*v1.Restore, error)
}

type RestoreIndexer func(obj *v1.Restore) ([]string, error)

type restoreController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewRestoreController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) RestoreController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &restoreController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromRestoreHandlerToHandler(sync RestoreHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Restore
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Restore))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *restoreController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Restore))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateRestoreDeepCopyOnChange(client RestoreClient, obj *v1.Restore, handler func(obj *v1.Restore) (*v1.Restore, error)) (*v1.Restore, error) {
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

func (c *restoreController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *restoreController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *restoreController) OnChange(ctx context.Context, name string, sync RestoreHandler) {
	c.AddGenericHandler(ctx, name, FromRestoreHandlerToHandler(sync))
}

func (c *restoreController) OnRemove(ctx context.Context, name string, sync RestoreHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromRestoreHandlerToHandler(sync)))
}

func (c *restoreController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *restoreController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *restoreController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *restoreController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *restoreController) Cache() RestoreCache {
	return &restoreCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *restoreController) Create(obj *v1.Restore) (*v1.Restore, error) {
	result := &v1.Restore{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *restoreController) Update(obj *v1.Restore) (*v1.Restore, error) {
	result := &v1.Restore{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *restoreController) UpdateStatus(obj *v1.Restore) (*v1.Restore, error) {
	result := &v1.Restore{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *restoreController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *restoreController) Get(namespace, name string, options metav1.GetOptions) (*v1.Restore, error) {
	result := &v1.Restore{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *restoreController) List(namespace string, opts metav1.ListOptions) (*v1.RestoreList, error) {
	result := &v1.RestoreList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *restoreController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *restoreController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.Restore, error) {
	result := &v1.Restore{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type restoreCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *restoreCache) Get(namespace, name string) (*v1.Restore, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.Restore), nil
}

func (c *restoreCache) List(namespace string, selector labels.Selector) (ret []*v1.Restore, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Restore))
	})

	return ret, err
}

func (c *restoreCache) AddIndexer(indexName string, indexer RestoreIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Restore))
		},
	}))
}

func (c *restoreCache) GetByIndex(indexName, key string) (result []*v1.Restore, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.Restore, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.Restore))
	}
	return result, nil
}

type RestoreStatusHandler func(obj *v1.Restore, status v1.RestoreStatus) (v1.RestoreStatus, error)

type RestoreGeneratingHandler func(obj *v1.Restore, status v1.RestoreStatus) ([]runtime.Object, v1.RestoreStatus, error)

func RegisterRestoreStatusHandler(ctx context.Context, controller RestoreController, condition condition.Cond, name string, handler RestoreStatusHandler) {
	statusHandler := &restoreStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromRestoreHandlerToHandler(statusHandler.sync))
}

func RegisterRestoreGeneratingHandler(ctx context.Context, controller RestoreController, apply apply.Apply,
	condition condition.Cond, name string, handler RestoreGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &restoreGeneratingHandler{
		RestoreGeneratingHandler: handler,
		apply:                    apply,
		name:                     name,
		gvk:                      controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterRestoreStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type restoreStatusHandler struct {
	client    RestoreClient
	condition condition.Cond
	handler   RestoreStatusHandler
}

func (a *restoreStatusHandler) sync(key string, obj *v1.Restore) (*v1.Restore, error) {
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
		var newErr error
		obj.Status = newStatus
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type restoreGeneratingHandler struct {
	RestoreGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *restoreGeneratingHandler) Remove(key string, obj *v1.Restore) (*v1.Restore, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.Restore{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *restoreGeneratingHandler) Handle(obj *v1.Restore, status v1.RestoreStatus) (v1.RestoreStatus, error) {
	objs, newStatus, err := a.RestoreGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
