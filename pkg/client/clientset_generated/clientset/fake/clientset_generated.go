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

package fake

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
	clientset "k8s.io/kops/pkg/client/clientset_generated/clientset"
	kopsinternalversion "k8s.io/kops/pkg/client/clientset_generated/clientset/typed/kops/internalversion"
	fakekopsinternalversion "k8s.io/kops/pkg/client/clientset_generated/clientset/typed/kops/internalversion/fake"
	kopsv1alpha1 "k8s.io/kops/pkg/client/clientset_generated/clientset/typed/kops/v1alpha1"
	fakekopsv1alpha1 "k8s.io/kops/pkg/client/clientset_generated/clientset/typed/kops/v1alpha1/fake"
	kopsv1alpha2 "k8s.io/kops/pkg/client/clientset_generated/clientset/typed/kops/v1alpha2"
	fakekopsv1alpha2 "k8s.io/kops/pkg/client/clientset_generated/clientset/typed/kops/v1alpha2/fake"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	fakePtr := testing.Fake{}
	fakePtr.AddReactor("*", "*", testing.ObjectReaction(o))
	fakePtr.AddWatchReactor("*", testing.DefaultWatchReactor(watch.NewFake(), nil))

	return &Clientset{fakePtr, &fakediscovery.FakeDiscovery{Fake: &fakePtr}}
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

var _ clientset.Interface = &Clientset{}

// Kops retrieves the KopsClient
func (c *Clientset) Kops() kopsinternalversion.KopsInterface {
	return &fakekopsinternalversion.FakeKops{Fake: &c.Fake}
}

// KopsV1alpha1 retrieves the KopsV1alpha1Client
func (c *Clientset) KopsV1alpha1() kopsv1alpha1.KopsV1alpha1Interface {
	return &fakekopsv1alpha1.FakeKopsV1alpha1{Fake: &c.Fake}
}

// KopsV1alpha2 retrieves the KopsV1alpha2Client
func (c *Clientset) KopsV1alpha2() kopsv1alpha2.KopsV1alpha2Interface {
	return &fakekopsv1alpha2.FakeKopsV1alpha2{Fake: &c.Fake}
}
