// Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package podzoneaffinity_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gardener/gardener/pkg/resourcemanager/webhook/podzoneaffinity"

	"github.com/gardener/gardener/pkg/logger"
	"github.com/gardener/gardener/pkg/operation/botanist/component/resourcemanager"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func TestPodZoneAffinity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PodZoneAffinity Integration Test Suite")
}

const testIDPrefix = "podzoneaffinity-webhook-test"

var (
	ctx = context.Background()
	log logr.Logger

	restConfig *rest.Config
	testEnv    *envtest.Environment
	testClient client.Client
)

var _ = BeforeSuite(func() {
	logf.SetLogger(logger.MustNewZapLogger(logger.DebugLevel, logger.FormatJSON, zap.WriteTo(GinkgoWriter)))
	log = logf.Log.WithName(testIDPrefix)

	By("starting test environment")
	testEnv = &envtest.Environment{
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			MutatingWebhooks: getMutatingWebhookConfigurations(),
		},
	}

	var err error
	restConfig, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(restConfig).NotTo(BeNil())

	DeferCleanup(func() {
		By("stopping test environment")
		Expect(testEnv.Stop()).To(Succeed())
	})

	By("creating test client")
	testClient, err = client.New(restConfig, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())

	By("setting up manager")
	mgr, err := manager.New(restConfig, manager.Options{
		Port:               testEnv.WebhookInstallOptions.LocalServingPort,
		Host:               testEnv.WebhookInstallOptions.LocalServingHost,
		CertDir:            testEnv.WebhookInstallOptions.LocalServingCertDir,
		MetricsBindAddress: "0",
		ClientDisableCacheFor: []client.Object{
			// Disable cache for namespaces so that changes applied by tests are seen immediately.
			&corev1.Namespace{},
		},
	})
	Expect(err).NotTo(HaveOccurred())

	By("registering webhook")
	Expect(podzoneaffinity.AddToManager(mgr)).To(Succeed())

	By("starting manager")
	mgrContext, mgrCancel := context.WithCancel(ctx)

	go func() {
		defer GinkgoRecover()
		Expect(mgr.Start(mgrContext)).To(Succeed())
	}()

	// Wait for the webhook server to start
	Eventually(func() error {
		checker := mgr.GetWebhookServer().StartedChecker()
		return checker(&http.Request{})
	}).Should(BeNil())

	DeferCleanup(func() {
		By("stopping manager")
		mgrCancel()
	})
})

func getMutatingWebhookConfigurations() []*admissionregistrationv1.MutatingWebhookConfiguration {
	webhookConfig := []*admissionregistrationv1.MutatingWebhookConfiguration{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: admissionregistrationv1.SchemeGroupVersion.String(),
				Kind:       "MutatingWebhookConfiguration",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "gardener-resource-manager",
			},
			Webhooks: []admissionregistrationv1.MutatingWebhook{
				resourcemanager.GetPodZoneAffinityMutatingWebhook(nil, func(_ *corev1.Secret, path string) admissionregistrationv1.WebhookClientConfig {
					return admissionregistrationv1.WebhookClientConfig{
						Service: &admissionregistrationv1.ServiceReference{
							Path: &path,
						},
					}
				}),
			},
		},
	}

	return webhookConfig
}