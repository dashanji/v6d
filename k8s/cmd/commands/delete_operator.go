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
package commands

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// deleteOperatorCmd deletes the vineyard operator on kubernetes
var deleteOperatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Delete the vineyard operator on kubernetes",
	Long: `Delete the vineyard operator on kubernetes. 
For example:

# delete the default vineyard operator in the vineyard-system namespace
vineyardctl delete operator

# delete the specific version of vineyard operator in the vineyard-system namespace
vineyardctl -n vineyard-system -k /home/gsbot/.kube/config delete operator -v 0.12.2

# delete the vineyard operator from local kustomize dir in the vineyard-system namespace
vineyardctl -n vineyard-system -k /home/gsbot/.kube/config delete operator --local ../config/default`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ValidateNoArgs("delete operator", args); err != nil {
			log.Fatal("failed to validate delete operator args and flags: ", err)
		}

		kubeClient, err := getKubeClient()
		if err != nil {
			log.Fatal("failed to get kubeclient: ", err)
		}

		operatorManifests, err := buildKustomizeDir(getKustomizeDir())
		if err != nil {
			log.Fatal("failed to build kustomize dir", err)
		}

		if err := deleteManifests(kubeClient, []byte(operatorManifests), GetDefaultVineyardNamespace()); err != nil {
			log.Fatal("failed to delete operator manifests: ", err)
		}

		waitOperatorDeleted(kubeClient)

		log.Println("Vineyard Operator is deleted.")
	},
}

// wait for the vineyard operator to be deleted
func waitOperatorDeleted(c client.Client) error {
	return wait.PollImmediate(1*time.Second, 300*time.Second, func() (bool, error) {
		deployment := &appsv1.Deployment{}
		err := c.Get(context.TODO(), types.NamespacedName{Name: "vineyard-controller-manager",
			Namespace: GetDefaultVineyardNamespace()}, deployment)
		if apierrors.IsNotFound(err) {
			return true, nil
		}
		return false, nil
	})
}

func NewDeleteOperatorCmd() *cobra.Command {
	return deleteOperatorCmd
}

func init() {
	deleteOperatorCmd.Flags().StringVarP(&OperatorVersion, "version", "v", "dev", "the version of kustomize dir from github repo")
	deleteOperatorCmd.Flags().StringVarP(&KustomzieDir, "local", "l", "", "the local kustomize dir")
}
