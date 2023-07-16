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
package sidecar

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	//"sigs.k8s.io/kustomize/kyaml/ext"

	//"github.com/stretchr/testify/assert"

	"github.com/v6d-io/v6d/k8s/apis/k8s/v1alpha1"
	"github.com/v6d-io/v6d/k8s/cmd/commands/flags"
	"github.com/v6d-io/v6d/k8s/cmd/commands/util"
	"github.com/v6d-io/v6d/k8s/controllers/k8s"
	"github.com/v6d-io/v6d/k8s/pkg/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidateFormat(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Valid YAML format",
			args:    args{format: YAMLFormat},
			wantErr: false,
		},
		{
			name:    "Valid JSON format",
			args:    args{format: JSONFormat},
			wantErr: false,
		},
		{
			name:    "Invalid format",
			args:    args{format: "invalid"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateFormat(tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("validateFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetWorkloadResource(t *testing.T) {
	tests := []struct {
		name           string
		workloadYAML   string
		workloadJSON   string
		expectedResult string
		expectedError  error
	}{
		{
			name:         "Valid YAML file",
			workloadYAML: "workload.yaml",
			workloadJSON: "",
			expectedResult: "apiVersion: apps/v1\n" +
				"kind: Deployment\n" +
				"metadata:\n" +
				"  name: nginx-deployment\n" +
				"  namespace: vineyard-system\n" +
				"spec:\n" +
				"  selector:\n" +
				"    matchLabels:\n" +
				"      app: nginx\n" +
				"  template:\n" +
				"    metadata:\n" +
				"      labels:\n" +
				"        app: nginx\n" +
				"    spec:\n" +
				"      containers:\n" +
				"      - name: nginx\n" +
				"        image: nginx:1.14.2\n" +
				"        ports:\n" +
				"        - containerPort: 80\n",

			expectedError: nil,
		},
		{
			name:           "Valid workload resource",
			workloadYAML:   "",
			workloadJSON:   `{"apiVersion":"apps/v1","kind":"Deployment"}`,
			expectedResult: "apiVersion: apps/v1\n" + "kind: Deployment\n",
			expectedError:  nil,
		},
		{
			name:           "Both workload yaml and workload resource specified",
			workloadYAML:   "workload.yaml",
			workloadJSON:   `{"apiVersion":"apps/v1","kind":"Deployment"}`,
			expectedResult: "",
			expectedError:  errors.New("cannot specify both workload resource and workload yaml"),
		},
	}

	for a, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags.WorkloadYaml = tt.workloadYAML
			flags.WorkloadResource = tt.workloadJSON
			result, err := getWorkloadResource()

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, but got no error", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if result != tt.expectedResult {
					fmt.Println(a)
					fmt.Println(result)
					fmt.Println(tt.expectedResult)
					t.Errorf("Expected result:\n%s\n\nBut got:\n%s", tt.expectedResult, result)

				}
			}
		})
	}
}

func TestGetWorkloadObj(t *testing.T) {
	workload := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
spec:
  template:
    spec:
      containers:
      - name: test-container
        image: nginx:latest
`
	expectedName := "test-deployment"
	expectedImage := "nginx:latest"

	obj, err := GetWorkloadObj(workload)
	assert.NoError(t, err)
	assert.NotNil(t, obj)

	name, found, err := unstructured.NestedString(obj.Object, "metadata", "name")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, expectedName, name)

	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "template", "spec", "containers")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Len(t, containers, 1)

	container := containers[0].(map[string]interface{})
	image, found, err := unstructured.NestedString(container, "image")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, expectedImage, image)
}

func TestGetManifestFromTemplate(t *testing.T) {
	flags.WorkloadYaml = "workload.yaml"
	resource, err := getWorkloadResource()
	if err != nil {
		log.Fatal(err, "failed to get the workload resource")
	}
	//sidecar, err := buildSidecar("vineyard-system")
	//fmt.Println(sidecar)
	//workloadObj, err := GetWorkloadObj(resource)
	//fmt.Println(workloadObj)
	_, err = GetManifestFromTemplate(resource)

	require.Error(t, err)

}

func TestParseManifestsAsYAML(t *testing.T) {
	om := OutputManifests{
		EtcdPod: []string{
			"etcd-pod-manifest-1",
			"etcd-pod-manifest-2",
		},
		EtcdInternalService: []string{
			"etcd-internal-service-manifest-1",
		},
		EtcdService: "etcd-service-manifest",
		RPCService:  "rpc-service-manifest",
		Workload:    "workload-manifest",
	}

	expectedYAML := []string{
		"etcd-pod-manifest-1\n",
		"etcd-pod-manifest-2\n",
		"etcd-internal-service-manifest-1\n",
		"etcd-service-manifest\n",
		"rpc-service-manifest\n",
		"workload-manifest\n",
	}

	// Mock ConvertToYaml function
	//util.ConvertToYaml = mockConvertToYaml

	yamlManifests, err := parseManifestsAsYAML(om)
	assert.NoError(t, err)
	assert.Equal(t, expectedYAML, yamlManifests)
}

func TestDeployDuringInjection(t *testing.T) {
	flags.KubeConfig = "/home/zhuyi/.kube/config"
	om := OutputManifests{
		Workload:    `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"creationTimestamp":null,"name":"nginx-deployment","namespace":"vineyard-system","ownerReferences":[]},"spec":{"selector":{"matchLabels":{"app":"nginx"}},"strategy":{},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"nginx","app.vineyard.io/name":"vineyard-sidecar"}},"spec":{"containers":[{"command":null,"image":"nginx:1.14.2","name":"nginx","ports":[{"containerPort":80}],"resources":{},"volumeMounts":[{"mountPath":"/var/run","name":"vineyard-socket"}]},{"command":["/bin/bash","-c","/usr/bin/wait-for-it.sh -t 60 vineyard-sidecar-etcd-service.vineyard-system.svc.cluster.local:2379; sleep 1; /usr/local/bin/vineyardd --sync_crds true --socket /var/run/vineyard.sock --size 256Mi --stream_threshold 80 --etcd_cmd etcd --etcd_prefix /vineyard --etcd_endpoint http://vineyard-sidecar-etcd-service:2379\n"],"env":[{"name":"VINEYARDD_UID","value":null},{"name":"VINEYARDD_NAME","value":"vineyard-sidecar"},{"name":"VINEYARDD_NAMESPACE","value":"vineyard-system"}],"image":"vineyardcloudnative/vineyardd:latest","imagePullPolicy":"IfNotPresent","name":"vineyard-sidecar","ports":[{"containerPort":9600,"name":"vineyard-rpc","protocol":"TCP"}],"resources":{"limits":null,"requests":null},"volumeMounts":[{"mountPath":"/var/run","name":"vineyard-socket"}]}],"volumes":[{"emptyDir":{},"name":"vineyard-socket"}]}}},"status":{}`,
		RPCService:  `{"apiVersion":"v1","kind":"Service","metadata":{"labels":{"app.vineyard.io/name":"vineyard-sidecar"},"name":"vineyard-sidecar-rpc","namespace":"vineyard-system","ownerReferences":[]},"spec":{"ports":[{"name":"vineyard-rpc","port":9600,"protocol":"TCP"}],"selector":{"app.vineyard.io/name":"vineyard-sidecar","app.vineyard.io/role":"vineyardd"},"type":"ClusterIP"}}`,
		EtcdService: `{"apiVersion":"v1","kind":"Service","metadata":{"name":"vineyard-sidecar-etcd-service","namespace":"vineyard-system","ownerReferences":[]},"spec":{"ports":[{"name":"vineyard-sidecar-etcd-for-vineyard-port","port":2379,"protocol":"TCP","targetPort":2379}],"selector":{"app.vineyard.io/name":"vineyard-sidecar","app.vineyard.io/role":"etcd"}}}`,
		EtcdInternalService: []string{
			`{"apiVersion":"v1","kind":"Service","metadata":{"labels":{"etcd_node":"vineyard-sidecar-etcd-0"},"name":"vineyard-sidecar-etcd-0","namespace":"vineyard-system","ownerReferences":[]},"spec":{"ports":[{"name":"client","port":2379,"protocol":"TCP","targetPort":2379},{"name":"server","port":2380,"protocol":"TCP","targetPort":2380}],"selector":{"app.vineyard.io/role":"etcd","etcd_node":"vineyard-sidecar-etcd-0"}}}`,
		},
		EtcdPod: []string{
			`{"apiVersion":"v1","kind":"Pod","metadata":{"labels":{"app.vineyard.io/name":"vineyard-sidecar","app.vineyard.io/role":"etcd","etcd_node":"vineyard-sidecar-etcd-0"},"name":"vineyard-sidecar-etcd-0","namespace":"vineyard-system","ownerReferences":[]},"spec":{"containers":[{"command":["etcd","--name","vineyard-sidecar-etcd-0","--initial-advertise-peer-urls","http://vineyard-sidecar-etcd-0:2380","--advertise-client-urls","http://vineyard-sidecar-etcd-0:2379","--listen-peer-urls","http://0.0.0.0:2380","--listen-client-urls","http://0.0.0.0:2379","--initial-cluster","vineyard-sidecar-etcd-0=http://vineyard-sidecar-etcd-0:2380","--initial-cluster-state","new"],"image":"vineyardcloudnative/vineyardd:latest","name":"etcd","ports":[{"containerPort":2379,"name":"client","protocol":"TCP"},{"containerPort":2380,"name":"server","protocol":"TCP"}]}],"restartPolicy":"Always"}}`,
		},
	}

	err := deployDuringInjection(&om)
	assert.NoError(t, err)

}

func TestOutputInjectedResult(t *testing.T) {
	om := OutputManifests{
		EtcdPod: []string{
			"etcd-pod-manifest-1",
			"etcd-pod-manifest-2",
		},
		EtcdInternalService: []string{
			"etcd-internal-service-manifest-1",
		},
		EtcdService: "etcd-service-manifest",
		RPCService:  "rpc-service-manifest",
		Workload:    "workload-manifest",
	}

	// Capture the output
	var buf bytes.Buffer
	//output = &buf

	// Test with JSON output format
	flags.OutputFormat = JSONFormat
	err := outputInjectedResult(om)
	assert.NoError(t, err)
	assert.Equal(t, "", buf.String())

	// Test with YAML output format
	buf.Reset()
	flags.OutputFormat = YAMLFormat
	err = outputInjectedResult(om)
	assert.NoError(t, err)
	assert.Equal(t, "", buf.String())
}

func TestBuildSidecar(t *testing.T) {
	namespace := "vineyard-system"
	// Set the flags and option

	sidecar, err := buildSidecar(namespace)
	assert.NoError(t, err)
	assert.NotNil(t, sidecar)
	assert.Equal(t, flags.SidecarName, sidecar.ObjectMeta.Name)
	assert.Equal(t, namespace, sidecar.ObjectMeta.Namespace)
	assert.Len(t, sidecar.Spec.Vineyard.Env, 0)

}

func TestInjectSidecarConfig(t *testing.T) {
	// Create a sidecar object
	sidecar := &v1alpha1.Sidecar{
		TypeMeta: metav1.TypeMeta{Kind: "", APIVersion: ""},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "vineyard-sidecar",
			GenerateName:               "",
			Namespace:                  "vineyard-system",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            []metav1.OwnerReference{},
			Finalizers:                 []string{},
			ZZZ_DeprecatedClusterName:  "",
			ManagedFields:              []metav1.ManagedFieldsEntry{},
		},
		Spec: v1alpha1.SidecarSpec{
			Selector:     "",
			Replicas:     1,
			EtcdReplicas: 0,
			Vineyard: v1alpha1.VineyardConfig{
				Image:           "vineyardcloudnative/vineyardd:latest",
				ImagePullPolicy: "IfNotPresent",
				SyncCRDs:        true,
				Socket:          "/var/run/vineyard-kubernetes/{{.Namespace}}/{{.Name}}",
				Size:            "256Mi",
				ReserveMemory:   false,
				StreamThreshold: 80,
				Spill: v1alpha1.SpillConfig{
					Name:                      "",
					Path:                      "",
					SpillLowerRate:            "0.3",
					SpillUpperRate:            "0.8",
					PersistentVolumeSpec:      corev1.PersistentVolumeSpec{},
					PersistentVolumeClaimSpec: corev1.PersistentVolumeClaimSpec{},
				},
				Env:    []corev1.EnvVar{},
				Memory: "",
				CPU:    "",
			},
			Metric: v1alpha1.MetricConfig{
				Enable:          false,
				Image:           "vineyardcloudnative/vineyard-grok-exporter:latest",
				ImagePullPolicy: "IfNotPresent",
			},
			Volume: v1alpha1.VolumeConfig{},
			Service: v1alpha1.ServiceConfig{
				Type: "clusterIP",
				Port: 9600,
			},
		},
		Status: v1alpha1.SidecarStatus{
			Current: 0,
		},
	}

	// Create a workload object
	workloadObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Deployment",
			"apiVersion": "apps/v1",
			"metadata": map[string]interface{}{
				"name":              "nginx-deployment",
				"namespace":         "vineyard-system",
				"creationTimestamp": nil,
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "nginx",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "nginx",
						},
						"creationTimestamp": nil,
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name": "nginx",
								"ports": []interface{}{
									map[string]interface{}{
										"containerPort": int64(80),
									},
								},
								"resources": map[string]interface{}{},
								"image":     "nginx:1.14.2",
							},
						},
					},
				},
				"strategy": map[string]interface{}{},
			},
			"status": map[string]interface{}{},
		},
	}

	// Create a sidecar object to inject
	var etcdConfig k8s.EtcdConfig
	tmplFunc := map[string]interface{}{
		"getEtcdConfig": func() k8s.EtcdConfig {
			return etcdConfig
		},
	}
	obj, err := util.RenderManifestAsObj(string("sidecar/injection-template.yaml"), sidecar, tmplFunc)

	err = InjectSidecarConfig(sidecar, workloadObj, obj)
	assert.NoError(t, err)
}

func TestNewInjectCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			want: injectCmd, // 指定预期的 *cobra.Command 值
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInjectCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInjectCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
