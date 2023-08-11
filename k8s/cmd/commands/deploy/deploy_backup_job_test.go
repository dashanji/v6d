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
package deploy

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/v6d-io/v6d/k8s/cmd/commands/flags"
	"github.com/v6d-io/v6d/k8s/cmd/commands/util"
	"sigs.k8s.io/controller-runtime/pkg/client"

	//"k8s.io/apimachinery/pkg/runtime/schema"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	//"k8s.io/client-go/kubernetes/fake"
	//"github.com/stretchr/testify/assert"
	//"k8s.io/api/apps/v1"
)

func TestDeployBackupJobCmd(t *testing.T) {
	// deploy a vineyardd for later backup operation
	flags.KubeConfig = os.Getenv("HOME") + "/.kube/config"
	//flags.KubeConfig = "/tmp/e2e-k8s.config"
	flags.Namespace = "vineyard-system"
	flags.VineyarddOpts.Replicas = 3
	flags.VineyarddOpts.EtcdReplicas = 1
	flags.VineyarddOpts.Vineyard.Image = "vineyardcloudnative/vineyardd:alpine-latest"
	flags.VineyarddOpts.Vineyard.CPU = ""
	flags.VineyarddOpts.Vineyard.Memory = ""
	flags.VineyarddOpts.Service.Port = 9600
	flags.VineyarddOpts.Service.Type = "ClusterIP"
	flags.VineyarddOpts.Volume.PvcName = ""
	flags.VineyarddOpts.Vineyard.Size = "256Mi"
	deployVineyardDeploymentCmd := NewDeployVineyardDeploymentCmd()
	deployVineyardDeploymentCmd.Run(deployVineyardDeploymentCmd, []string{})

	//backup operation
	flags.BackupOpts.BackupPath = "/var/vineyard/dump"
	//flags.PVCName = "pvc-for-backup-and-recover-demo"

	flags.VineyardDeploymentName = "vineyardd-sample"
	flags.VineyardDeploymentNamespace = "vineyard-system"
	flags.BackupPVandPVC = `{
		"pv-spec": {
		  "capacity": {
			"storage": "1Gi"
		  },
		  "accessModes": [
			"ReadWriteOnce"
		  ],
		  "storageClassName": "manual",
		  "hostPath": {
			"path": "/var/vineyard/dump"
		  }
		},
		"pvc-spec": {
		  "storageClassName": "manual",
		  "accessModes": [
			"ReadWriteOnce"
		  ],
		  "resources": {
			"requests": {
			"storage": "1Gi"
            }
          },
		  "volumeName" : "backup-path",
          "selector": {
			"matchLabels": {
				"app.kubernetes.io/name" : "vineyard-backup"
			}
		  }
		}
	}
    `
	c := util.KubernetesClient()
	got, _ := getBackupObjectsFromTemplate(c, []string{})
	for i := range got {
		fmt.Println(*got[i])
	}
	deployBackupJobCmd.Run(deployBackupJobCmd, []string{})

	if util.Wait(func() (bool, error) {
		jobName := flags.BackupName
		name := client.ObjectKey{Name: jobName, Namespace: flags.Namespace}
		job := batchv1.Job{}
		if err := c.Get(context.TODO(), name, &job); err != nil {
			return false, err
		}
		if job.Status.Succeeded == *job.Spec.Parallelism {
			return true, nil
		}
		return false, nil
	}) != nil {
		t.Errorf("backup job can not be deployed successfully")
	}
}

func Test_getBackupObjectsFromTemplate(t *testing.T) {
	// set the flags
	flags.KubeConfig = os.Getenv("HOME") + "/.kube/config"
	//flags.KubeConfig = "/tmp/e2e-k8s.config"
	flags.PVCName = "pvc-for-backup-and-recover-demo"
	flags.VineyarddOpts.Volume.PvcName = "pvc-for-backup-and-recover-demo"
	flags.Namespace = "vineyard-system"
	flags.VineyardDeploymentName = "vineyardd-sample"
	flags.VineyardDeploymentNamespace = "vineyard-system"
	c := util.KubernetesClient()

	type args struct {
		c    client.Client
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*unstructured.Unstructured
		wantErr bool
	}{
		// Add test cases.
		{
			name: "Test case ",
			args: args{
				c:    c,
				args: []string{},
			},
			want: []*unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "batch/v1",
						"kind":       "Job",
						"metadata": map[string]interface{}{
							"name":      "vineyard-backup",
							"namespace": "vineyard-system",
						},
						"spec": map[string]interface{}{
							"parallelism": int64(3),
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"labels": map[string]interface{}{
										"app.kubernetes.io/name": "vineyard-backup",
									},
								},
								"spec": map[string]interface{}{
									"affinity": map[string]interface{}{
										"podAffinity": map[string]interface{}{
											"requiredDuringSchedulingIgnoredDuringExecution": []interface{}{
												map[string]interface{}{
													"labelSelector": map[string]interface{}{
														"matchExpressions": []interface{}{
															map[string]interface{}{
																"key":      "app.kubernetes.io/instance",
																"operator": "In",
																"values": []interface{}{
																	"vineyard-system-vineyardd-sample",
																},
															},
														},
													},
													"topologyKey": "kubernetes.io/hostname",
												},
											},
										},
										"podAntiAffinity": map[string]interface{}{
											"requiredDuringSchedulingIgnoredDuringExecution": []interface{}{
												map[string]interface{}{
													"labelSelector": map[string]interface{}{
														"matchExpressions": []interface{}{
															map[string]interface{}{
																"key":      "app.kubernetes.io/name",
																"operator": "In",
																"values": []interface{}{
																	"vineyard-backup",
																},
															},
														},
													},
													"topologyKey": "kubernetes.io/hostname",
												},
											},
										},
									},
									"containers": []interface{}{
										map[string]interface{}{
											"env": []interface{}{

												map[string]interface{}{
													"name":  "BACKUP_PATH",
													"value": "/var/vineyard/dump",
												},
												map[string]interface{}{
													"name":  "ENDPOINT",
													"value": "vineyardd-sample-rpc.vineyard-system",
												},
												map[string]interface{}{
													"name":  "SELECTOR",
													"value": "vineyard-backup",
												},
												map[string]interface{}{
													"name":  "ALLINSTANCES",
													"value": "3",
												},
												map[string]interface{}{
													"name": "POD_NAME",
													"valueFrom": map[string]interface{}{
														"fieldRef": map[string]interface{}{
															"fieldPath": "metadata.name",
														},
													},
												},
												map[string]interface{}{
													"name": "POD_NAMESPACE",
													"valueFrom": map[string]interface{}{
														"fieldRef": map[string]interface{}{
															"fieldPath": "metadata.namespace",
														},
													},
												},
											},
											"image":           "ghcr.io/v6d-io/v6d/backup-job",
											"imagePullPolicy": "IfNotPresent",
											"name":            "engine",
											"volumeMounts": []interface{}{
												map[string]interface{}{
													"mountPath": "/var/run",
													"name":      "vineyard-sock",
												},
												map[string]interface{}{
													"mountPath": "/var/vineyard/dump",
													"name":      "backup-path",
												},
											},
										},
									},
									"restartPolicy": "Never",
									"volumes": []interface{}{
										map[string]interface{}{
											"hostPath": map[string]interface{}{
												"path": "/var/run/vineyard-kubernetes/vineyard-system/vineyardd-sample",
											},
											"name": "vineyard-sock",
										},
										map[string]interface{}{
											"name": "backup-path",
											"persistentVolumeClaim": map[string]interface{}{
												"claimName": "pvc-for-backup-and-recover-demo",
											},
										},
									},
								},
							},
							"ttlSecondsAfterFinished": int64(80),
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "rbac.authorization.k8s.io/v1",
						"kind":       "RoleBinding",
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"app.kubernetes.io/name": "backup",
							},
							"name":      "vineyard-backup",
							"namespace": "vineyard-system",
						},
						"roleRef": map[string]interface{}{
							"apiGroup": "rbac.authorization.k8s.io",
							"kind":     "Role",
							"name":     "vineyard-backup",
						},
						"subjects": []interface{}{
							map[string]interface{}{
								"kind":      "ServiceAccount",
								"name":      "default",
								"namespace": "vineyard-system",
							},
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "rbac.authorization.k8s.io/v1",
						"kind":       "Role",
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"app.kubernetes.io/instance": "backup",
							},
							"name":      "vineyard-backup",
							"namespace": "vineyard-system",
						},
						"rules": []interface{}{
							map[string]interface{}{
								"apiGroups": []interface{}{
									"",
								},
								"resources": []interface{}{
									"pods",
									"pods/log",
								},
								"verbs": []interface{}{
									"get",
									"list",
								},
							},
							map[string]interface{}{
								"apiGroups": []interface{}{
									"",
								},
								"resources": []interface{}{
									"pods/exec",
								},
								"verbs": []interface{}{
									"create",
								},
							},
						},
						"subjects": []interface{}{
							map[string]interface{}{
								"kind":      "ServiceAccount",
								"name":      "default",
								"namespace": "vineyard-system",
							},
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "PersistentVolumeClaim",
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"app.kubernetes.io/name": "vineyard-backup",
							},
							"name":      "vineyard-backup",
							"namespace": "vineyard-system",
						},
						"spec": map[string]interface{}{
							"resources": nil,
							"selector": map[string]interface{}{
								"matchLabels": map[string]interface{}{
									"app.kubernetes.io/name": "vineyard-backup",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBackupObjectsFromTemplate(tt.args.c, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBackupObjectsFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for i := range got {
				if !reflect.DeepEqual(*got[i], *(tt.want)[i]) {
					fmt.Println(*got[i])
					t.Errorf("getBackupObjectsFromTemplate() = %+v, want %+v", got, tt.want)

				}
			}

		})
	}
}