/*
* Copyright 2020-2023 Alibaba Group Holding Limited.

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
package util

import (
	"fmt"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	vineyardV1alpha1 "github.com/v6d-io/v6d/k8s/apis/k8s/v1alpha1"
	"github.com/v6d-io/v6d/k8s/cmd/commands/flags"
	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoScheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ClientgoScheme    = runtime.NewScheme()
	OperatorScheme    = runtime.NewScheme()
	CertManagerScheme = runtime.NewScheme()
)

// AddClientGoScheme add client-go scheme to CmdScheme
func AddClientGoScheme(scheme *runtime.Scheme) error {
	err := clientgoScheme.AddToScheme(scheme)
	if err != nil {
		return fmt.Errorf("failed to add client-go scheme to current scheme: %v", err)
	}
	return nil
}

// AddVineyardV1alpha1Scheme add vineyard v1alpha1 scheme to CmdScheme
func AddVineyardV1alpha1Scheme(scheme *runtime.Scheme) error {
	err := vineyardV1alpha1.AddToScheme(scheme)
	if err != nil {
		return fmt.Errorf("failed to add vineyard v1alpha1 scheme to current scheme: %v", err)
	}
	return nil
}

// AddApiExtensionsScheme add apiextensions scheme to CmdScheme
func AddApiExtensionsScheme(scheme *runtime.Scheme) error {
	err := apiv1.AddToScheme(scheme)
	if err != nil {
		return fmt.Errorf("failed to add apiextensions scheme to current scheme: %v", err)
	}
	return nil
}

// AddCertManagerScheme add cert-manager scheme to CmdScheme
func AddCertManagerScheme(scheme *runtime.Scheme) error {
	err := cmapi.AddToScheme(scheme)
	if err != nil {
		return fmt.Errorf("failed to add cert-manager scheme to current scheme: %v", err)
	}
	return nil
}

// GetClientgoScheme return the client-go scheme
func GetClientgoScheme() (*runtime.Scheme, error) {
	if err := AddClientGoScheme(ClientgoScheme); err != nil {
		return nil, err
	}
	return ClientgoScheme, nil
}

// GetOperatorScheme return the operator scheme
func GetOperatorScheme() (*runtime.Scheme, error) {
	if err := AddClientGoScheme(OperatorScheme); err != nil {
		return nil, err
	}
	if err := AddCertManagerScheme(OperatorScheme); err != nil {
		return nil, err
	}
	if err := AddVineyardV1alpha1Scheme(OperatorScheme); err != nil {
		return nil, err
	}
	if err := AddApiExtensionsScheme(OperatorScheme); err != nil {
		return nil, err
	}
	return OperatorScheme, nil
}

// GetCertManagerScheme return the cert-manager scheme
func GetCertManagerScheme() (*runtime.Scheme, error) {
	if err := AddClientGoScheme(CertManagerScheme); err != nil {
		return nil, err
	}
	if err := AddCertManagerScheme(CertManagerScheme); err != nil {
		return nil, err
	}
	if err := AddApiExtensionsScheme(CertManagerScheme); err != nil {
		return nil, err
	}
	return CertManagerScheme, nil
}

// GetKubeClient return the kubernetes client
func GetKubeClient(scheme *runtime.Scheme) (client.Client, error) {
	cfg := &rest.Config{}

	cfg, err := clientcmd.BuildConfigFromFlags("", flags.Kubeconfig)
	if err != nil {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	client, _ := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return client, nil
}