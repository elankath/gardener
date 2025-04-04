// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	corev1beta1 "github.com/gardener/gardener/pkg/client/core/clientset/versioned/typed/core/v1beta1"
	gentype "k8s.io/client-go/gentype"
)

// fakeControllerDeployments implements ControllerDeploymentInterface
type fakeControllerDeployments struct {
	*gentype.FakeClientWithList[*v1beta1.ControllerDeployment, *v1beta1.ControllerDeploymentList]
	Fake *FakeCoreV1beta1
}

func newFakeControllerDeployments(fake *FakeCoreV1beta1) corev1beta1.ControllerDeploymentInterface {
	return &fakeControllerDeployments{
		gentype.NewFakeClientWithList[*v1beta1.ControllerDeployment, *v1beta1.ControllerDeploymentList](
			fake.Fake,
			"",
			v1beta1.SchemeGroupVersion.WithResource("controllerdeployments"),
			v1beta1.SchemeGroupVersion.WithKind("ControllerDeployment"),
			func() *v1beta1.ControllerDeployment { return &v1beta1.ControllerDeployment{} },
			func() *v1beta1.ControllerDeploymentList { return &v1beta1.ControllerDeploymentList{} },
			func(dst, src *v1beta1.ControllerDeploymentList) { dst.ListMeta = src.ListMeta },
			func(list *v1beta1.ControllerDeploymentList) []*v1beta1.ControllerDeployment {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1beta1.ControllerDeploymentList, items []*v1beta1.ControllerDeployment) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
