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
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"fmt"

	"github.com/v6d-io/v6d/k8s/cmd/commands/flags"
	"github.com/v6d-io/v6d/k8s/cmd/commands/util"
)

func TestDeployVineyardDeploymentCmd(t *testing.T) {
	testReplicas := struct {
		name             string
		vineyardReplicas int
		etcdReplicas     int
	}{
		name:             "test replicas",
		vineyardReplicas: 1,
		etcdReplicas:     3,
	}
	t.Run(testReplicas.name, func(t *testing.T) {
		// set the flags
		flags.KubeConfig = "/tmp/e2e-k8s.config"
		flags.VineyarddOpts.Replicas = 1
		flags.VineyarddOpts.EtcdReplicas = 3
		deployVineyardDeploymentCmd.Run(deployVineyardDeploymentCmd, []string{})
		// get the replicas of etcd and vineyardd
		k8sclient := util.KubernetesClient()
		vineyardPods := corev1.PodList{}
		vineyarddOpts := []client.ListOption{
			client.InNamespace(flags.Namespace),
			client.MatchingLabels{
				"app.vineyard.io/name": flags.VineyarddName,
				"app.vineyard.io/role": "vineyardd",
			},
		}
		err := k8sclient.List(context.Background(), &vineyardPods, vineyarddOpts...)
		if err != nil {
			t.Errorf("list vineyardd pods error: %v", err)
		}

		etcdPod := corev1.PodList{}
		etcdOpts := []client.ListOption{
			client.InNamespace(flags.Namespace),
			client.MatchingLabels{
				"app.vineyard.io/name": flags.VineyarddName,
				"app.vineyard.io/role": "etcd",
			},
		}
		err = k8sclient.List(context.Background(), &etcdPod, etcdOpts...)
		if err != nil {
			t.Errorf("list etcd pods error: %v", err)
		}

		if len(vineyardPods.Items) != testReplicas.vineyardReplicas {
			t.Errorf("vineyardd replicas want: %d, got: %d", testReplicas.vineyardReplicas, len(vineyardPods.Items))
		}

		if len(etcdPod.Items) != testReplicas.etcdReplicas {
			t.Errorf("etcd replicas want: %d, got: %d", testReplicas.etcdReplicas, len(etcdPod.Items))
		}

	})
}

func TestGetVineyardDeploymentObjectsFromTemplate(t *testing.T) {
	tests := []struct {
		name    string
		want    []*unstructured.Unstructured
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test case",
			want: []*unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"app.vineyard.io/name": "vineyardd-sample",
								"app.vineyard.io/role": "etcd",
								"etcd_node":            "vineyardd-sample-etcd-0",
							},
							"name":      "vineyardd-sample-etcd-0",
							"namespace": "vineyard-system",
						},
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"command": []interface{}{
										"etcd",
										"--name",
										"vineyardd-sample-etcd-0",
										"--initial-advertise-peer-urls",
										"http://vineyardd-sample-etcd-0:2380",
										"--advertise-client-urls",
										"http://vineyardd-sample-etcd-0:2379",
										"--listen-peer-urls",
										"http://0.0.0.0:2380",
										"--listen-client-urls",
										"http://0.0.0.0:2379",
										"--initial-cluster",
										"vineyardd-sample-etcd-0=http://vineyardd-sample-etcd-0:2380",
										"--initial-cluster-state",
										"new",
									},
									"image": "vineyardcloudnative/vineyardd:latest",
									"name":  "etcd",
									"ports": []interface{}{
										map[string]interface{}{
											"containerPort": int64(2379),
											"name":          "client",
											"protocol":      "TCP",
										},
										map[string]interface{}{
											"containerPort": int64(2380),
											"name":          "server",
											"protocol":      "TCP",
										},
									},
								},
							},

							"restartPolicy": "Always",
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Service",
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"etcd_node": "vineyardd-sample-etcd-0",
							},
							"name":      "vineyardd-sample-etcd-0",
							"namespace": "vineyard-system",
						},
						"spec": map[string]interface{}{
							"ports": []interface{}{
								map[string]interface{}{
									"name":       "client",
									"port":       int64(2379),
									"protocol":   "TCP",
									"targetPort": int64(2379),
								},
								map[string]interface{}{
									"name":       "server",
									"port":       int64(2380),
									"protocol":   "TCP",
									"targetPort": int64(2380),
								},
							},
							"selector": map[string]interface{}{
								"app.vineyard.io/role": "etcd",
								"etcd_node":            "vineyardd-sample-etcd-0",
							},
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "apps/v1",
						"kind":       "Deployment",
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"app.kubernetes.io/component": "deployment",
								"app.kubernetes.io/instance":  "vineyard-system-vineyardd-sample",
								"app.vineyard.io/name":        "vineyardd-sample",
							},
							"name":      "vineyardd-sample",
							"namespace": "vineyard-system",
						},
						"spec": map[string]interface{}{
							"replicas": int64(3),
							"selector": map[string]interface{}{
								"matchLabels": map[string]interface{}{
									"app.kubernetes.io/instance": "vineyard-system-vineyardd-sample",
									"app.kubernetes.io/name":     "vineyardd-sample",
									"app.vineyard.io/name":       "vineyardd-sample",
								},
							},
							"template": map[string]interface{}{
								"metadata": map[string]interface{}{
									"annotations": map[string]interface{}{
										"kubectl.kubernetes.io/default-container":      "vineyardd",
										"kubectl.kubernetes.io/default-logs-container": "vineyardd",
									},
									"labels": map[string]interface{}{
										"app.kubernetes.io/component": "deployment",
										"app.kubernetes.io/instance":  "vineyard-system-vineyardd-sample",
										"app.kubernetes.io/name":      "vineyardd-sample",
										"app.vineyard.io/name":        "vineyardd-sample",
										"app.vineyard.io/role":        "vineyardd",
									},
								},
								"spec": map[string]interface{}{
									"affinity": map[string]interface{}{
										"podAntiAffinity": map[string]interface{}{
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
									},
									"containers": []interface{}{
										map[string]interface{}{
											"env": []interface{}{
												map[string]interface{}{
													"name":  "VINEYARDD_UID",
													"value": nil,
												},
												map[string]interface{}{
													"name":  "VINEYARDD_NAME",
													"value": "vineyardd-sample",
												},
												map[string]interface{}{
													"name":  "VINEYARDD_NAMESPACE",
													"value": "vineyard-system",
												},
												map[string]interface{}{
													"name": "MY_NODE_NAME",
													"valueFrom": map[string]interface{}{
														"fieldRef": map[string]interface{}{
															"fieldPath": "spec.nodeName",
														},
													},
												},
												map[string]interface{}{
													"name": "MY_POD_NAME",
													"valueFrom": map[string]interface{}{
														"fieldRef": map[string]interface{}{
															"fieldPath": "metadata.name",
														},
													},
												},
												map[string]interface{}{
													"name": "MY_POD_NAMESPACE",
													"valueFrom": map[string]interface{}{
														"fieldRef": map[string]interface{}{
															"fieldPath": "metadata.namespace",
														},
													},
												},
												map[string]interface{}{
													"name": "MY_UID",
													"valueFrom": map[string]interface{}{
														"fieldRef": map[string]interface{}{
															"fieldPath": "metadata.uid",
														},
													},
												},
												map[string]interface{}{
													"name": "MY_POD_IP",
													"valueFrom": map[string]interface{}{
														"fieldRef": map[string]interface{}{
															"fieldPath": "status.podIP",
														},
													},
												},
												map[string]interface{}{
													"name": "MY_HOST_NAME",
													"valueFrom": map[string]interface{}{
														"fieldRef": map[string]interface{}{
															"fieldPath": "status.podIP",
														},
													},
												},
												map[string]interface{}{
													"name": "USER",
													"valueFrom": map[string]interface{}{
														"fieldRef": map[string]interface{}{
															"fieldPath": "metadata.name",
														},
													},
												},
											},
											"readinessProbe": map[string]interface{}{
												"exec": map[string]interface{}{
													"command": []interface{}{
														"ls",
														"/var/run/vineyard.sock",
													},
												},
											},
											"resources": map[string]interface{}{
												"limits":   nil,
												"requests": nil,
											},
											"command": []interface{}{
												"/bin/bash",
												"-c",
												"/usr/bin/wait-for-it.sh -t 60 vineyardd-sample-etcd-service.vineyard-system.svc.cluster.local:2379; sleep 1; /usr/local/bin/vineyardd --sync_crds true --socket /var/run/vineyard.sock --size 256Mi --stream_threshold 80 --etcd_cmd etcd --etcd_prefix /vineyard --etcd_endpoint http://vineyardd-sample-etcd-service:2379\n",
											},
											"image":           "vineyardcloudnative/vineyardd:latest",
											"imagePullPolicy": "IfNotPresent",
											"livenessProbe": map[string]interface{}{
												"periodSeconds": int64(60),
												"tcpSocket": map[string]interface{}{
													"port": int64(9600),
												},
											},
											"name": "vineyardd",
											"ports": []interface{}{
												map[string]interface{}{
													"containerPort": int64(9600),
													"name":          "rpc",
													"protocol":      "TCP",
												},
											},
											"volumeMounts": []interface{}{
												map[string]interface{}{
													"mountPath": "/var/run",
													"name":      "vineyard-socket",
												},
												map[string]interface{}{
													"mountPath": "/dev/shm",
													"name":      "shm",
												},
												map[string]interface{}{
													"mountPath": "/var/log/vineyard",
													"name":      "log",
												},
											},
										},
									},
									"volumes": []interface{}{
										map[string]interface{}{
											"hostPath": map[string]interface{}{
												"path": "/var/run/vineyard-kubernetes/vineyard-system/vineyardd-sample",
											},
											"name": "vineyard-socket",
										},
										map[string]interface{}{
											"emptyDir": map[string]interface{}{
												"medium": "Memory",
											},
											"name": "shm",
										},
										map[string]interface{}{
											"emptyDir": map[string]interface{}{},
											"name":     "log",
										},
									},
								},
							},
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Service",
						"metadata": map[string]interface{}{
							"name":      "vineyardd-sample-etcd-service",
							"namespace": "vineyard-system",
						},
						"spec": map[string]interface{}{
							"ports": []interface{}{
								map[string]interface{}{
									"name":       "vineyardd-sample-etcd-for-vineyard-port",
									"port":       int64(2379),
									"protocol":   "TCP",
									"targetPort": int64(2379),
								},
							},
							"selector": map[string]interface{}{
								"app.vineyard.io/name": "vineyardd-sample",
								"app.vineyard.io/role": "etcd",
							},
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Service",
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"app.vineyard.io/name": "vineyardd-sample",
							},
							"name":      "vineyardd-sample-rpc",
							"namespace": "vineyard-system",
						},
						"spec": map[string]interface{}{
							"ports": []interface{}{
								map[string]interface{}{
									"name":     "vineyard-rpc",
									"port":     int64(9600),
									"protocol": "TCP",
								},
							},
							"selector": map[string]interface{}{
								"app.vineyard.io/name": "vineyardd-sample",
								"app.vineyard.io/role": "vineyardd",
							},
							"type": "ClusterIP",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags.Namespace = "vineyard-system"
			got, err := GetVineyardDeploymentObjectsFromTemplate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVineyardDeploymentObjectsFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(*got[i], *(tt.want)[i]) {
					fmt.Println(i)
					fmt.Println(*got[i])
					fmt.Println(*(tt.want)[i])
					t.Errorf("getBackupObjectsFromTemplate() = %+v, want %+v", got, tt.want)

				}
			}
		})
	}
}

func Test_applyVineyarddFromTemplate(t *testing.T) {
	flags.KubeConfig = "/home/zhuyi/.kube/config"
	c := util.KubernetesClient()

	type args struct {
		c client.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test case",
			args: args{
				c: c,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags.Namespace = "vineyard-system"
			if err := applyVineyarddFromTemplate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("applyVineyarddFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/*func Test_waitVineyardDeploymentReady(t *testing.T) {
	flags.KubeConfig = "/home/zhuyi/.kube/config"
	c := util.KubernetesClient()
	type args struct {
		c client.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test case",
			args: args{
				c: c,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags.Namespace = "vineyard-system"
			if err := waitVineyardDeploymentReady(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("waitVineyardDeploymentReady() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}*/
