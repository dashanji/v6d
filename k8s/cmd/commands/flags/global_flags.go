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
package flags

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const defaultNamespace = "vineyard-system"

var (
	// kubeconfig path
	KubeConfig string

	// Namespace for operation
	Namespace string

	// DumpUsage
	DumpUsage bool
)

// defaultKubeConfig return the default kubeconfig path
func defaultKubeConfig() string {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}
	return kubeconfig
}

// GetDefaultVineyardNamespace return the default vineyard namespace
func GetDefaultVineyardNamespace() string {
	// we don't use the default namespace for vineyard
	if Namespace == "default" {
		return defaultNamespace
	}
	return Namespace
}

func ApplyGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().
		StringVarP(&KubeConfig, "kubeconfig", "", defaultKubeConfig(), "kubeconfig path for the kubernetes cluster")
	cmd.PersistentFlags().
		StringVarP(&Namespace, "namespace", "n", defaultNamespace, "the namespace for operation")
	cmd.PersistentFlags().BoolVarP(&DumpUsage, "dump-usage", "j", false, "Dump usage in JSON")
}

func RemoveVersionFlag(f *pflag.FlagSet) {
	// Avoid conflict of customized "version" flag with the one registered
	// in the `verflag` package
	normalize := f.GetNormalizeFunc()
	f.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		if name == "version" /* the argument `f` should be pflag.CommandLine */ {
			return pflag.NormalizedName("x-version")
		}
		return pflag.NormalizedName("version")
	})
	// restore
	f.SetNormalizeFunc(normalize)

	// hidden it from the help/usage
	f.Lookup("x-version").Hidden = true
}
