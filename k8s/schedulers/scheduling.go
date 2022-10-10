/*
Copyright 2020 The Kubernetes Authors.
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

package schedulers

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"k8s.io/klog/v2"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/component-helpers/scheduling/corev1"

	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	v1alpha1 "github.com/v6d-io/v6d/k8s/apis/k8s/v1alpha1"
)

const (
	// Name is the name of the plugin used in Registry and configurations.
	Name = "Vineyard"
	// Timeout is the default timeout for the scheduler plugin.
	Timeout = 60
	// VineyardJobName is the pod group name
	VineyardJobName = "scheduling.k8s.v6d.io/job"
	// VineyardJobRequired is the object ids that required by this job
	VineyardJobRequired = "scheduling.k8s.v6d.io/required"
	// VineyardJobReplica is the replication of pods in this job.
	VineyardJobReplica = "scheduling.k8s.v6d.io/replica"
	// ControlPlaneLabel is the label of the control plane
	ControlPlaneLabel = "node-role.kubernetes.io/control-plane"
	// VineyardSystemNamespace is the default system namespace
	VineyardSystemNamespace = "vineyard-system"
	// VineyarddName is the name of the vineyardd
	VineyarddName = "scheduling.k8s.v6d.io/vineyardd"
)

// SchedulerState records the status of current scheduling
type SchedulerState struct {
	client.Client
	state map[string]map[string]string // { jobname: { pod: nodename }}
}

// Append records the action of appending a new pod in job to given node.
func (ss *SchedulerState) Append(job string, pod string, nodeName string) error {
	klog.V(5).Infof("assign job %v pod %v to node %v", job, pod, nodeName)
	if s, ok := ss.state[job]; ok {
		if _, ok := s[pod]; ok {
			return fmt.Errorf("The pod has already been scheduled")
		}
		s[pod] = nodeName
		return nil
	}
	ss.state[job] = make(map[string]string)
	ss.state[job][pod] = nodeName
	return nil
}

// Compute the placement of a pod in job, assuming the useable nodes, and based on the given objects pool.
// Use the
func (ss *SchedulerState) Compute(ctx context.Context, job string, replica int64, rank int64,
	workernodes []string, requires []string, nodeName string, pod *v1.Pod) (int64, error) {
	// if requires no vineyard object, the job can be deployed in any nodes.
	// use round-robin scheduling here
	num := len(workernodes)
	if len(requires) == 0 {
		if workernodes[int(rank)%num] == nodeName {
			klog.V(5).Infof("nodeName: %v, workernodes: %v, rank: %v", nodeName, workernodes, rank)
			return 100, nil
		} else {
			return 1, nil
		}
	}
	// if no replica, raise
	if replica == 0 {
		return 0, fmt.Errorf("No replica information in the job spec")
	}

	// if the pod needs to be injected with the assembly container
	// we must wait for the assembly container to be ready
	if value, ok := pod.Labels["assembly.v6d.io/enabled"]; ok && strings.ToLower(value) == "true" {
		return 0, nil
	}

	// accumulates all local required objects
	globalObjects, err := ss.getGlobalObjectsByID(ctx, requires)
	if err != nil {
		return 0, err
	}
	klog.V(5).Infof("job %v requires objects %v", job, globalObjects)
	localsigs := make([]string, 0)
	for _, globalObject := range globalObjects {
		localsigs = append(localsigs, globalObject.Spec.Members...)
	}
	localObjects, err := ss.getLocalObjectsBySignatures(ctx, localsigs)
	if err != nil {
		return 0, err
	}
	if len(localObjects) == 0 {
		return 0, fmt.Errorf("No local chunks found")
	}

	if err := ss.createConfigmapForID(ctx, requires, pod.GetNamespace(), localObjects, globalObjects, pod); err != nil {
		klog.V(5).Infof("can't create configmap for object ID %v", err)
	}

	klog.V(5).Infof("job %v requires local chunks %v", job, localObjects)

	locations := make(map[string][]string)
	for _, localObject := range localObjects {
		host := localObject.Spec.Hostname
		if _, ok := locations[host]; !ok {
			locations[host] = make([]string, 0)
		}
		locations[host] = append(locations[host], localObject.Spec.ObjectID)
	}

	// total frags
	totalfrags := int64(len(localObjects))
	// frags for per pod
	nchunks := totalfrags / replica
	if totalfrags%replica != 0 {
		nchunks++
	}

	// find the node
	nodes := make([]string, 0)
	for k := range locations {
		nodes = append(nodes, k)
	}
	sort.Strings(nodes)

	var cnt int64 = 0
	target := ""
	for _, node := range nodes {
		localfrags := int64(len(locations[node]))
		if cnt+localfrags >= (nchunks*rank + (nchunks+1)/2) {
			target = node
			break
		}
		cnt += localfrags
	}

	// make sure every pod will be deployed in a node
	if target == "" {
		if nodeName == nodes[0] {
			klog.V(5).Infof("Bint the pod to the node with the most locations, %v", nodes[0])
			return 100, nil
		}
		return 1, nil
	} else if target == nodeName {
		klog.V(5).Infof("target == nodeName")
		return 100, nil
	} else {
		return 1, nil
	}
}

func (ss *SchedulerState) getGlobalObjectsByID(ctx context.Context, jobNames []string) ([]*v1alpha1.GlobalObject, error) {
	requiredjobs := make(map[string]bool)
	for _, n := range jobNames {
		requiredjobs[n] = true
	}
	objects := []*v1alpha1.GlobalObject{}
	globalObjects := &v1alpha1.GlobalObjectList{}
	if err := ss.List(ctx, globalObjects); err != nil {
		klog.V(5).Infof("client.List failed to get global objects, error: %v", err)
		return nil, err
	}
	for _, obj := range globalObjects.Items {
		if jobname, exist := obj.Labels["job"]; exist && requiredjobs[jobname] {
			objects = append(objects, &obj)
		}
	}

	return objects, nil
}

func (ss *SchedulerState) getLocalObjectsBySignatures(ctx context.Context, signatures []string) ([]*v1alpha1.LocalObject, error) {
	objects := make([]*v1alpha1.LocalObject, 0)
	for _, sig := range signatures {
		localObjects := &v1alpha1.LocalObjectList{}
		if err := ss.List(ctx, localObjects, client.MatchingLabels{
			"k8s.v6d.io/signature": sig,
		}); err != nil {
			klog.V(5).Infof("client.List failed to get local objects, error: %v", err)
			return nil, err
		} else {
			for _, localObject := range localObjects.Items {
				objects = append(objects, &localObject)
			}
		}
	}

	return objects, nil
}

// Create a configmap for the object id and the nodes
func (ss *SchedulerState) createConfigmapForID(ctx context.Context, jobname []string, namespace string,
	localobjects []*v1alpha1.LocalObject, globalobjects []*v1alpha1.GlobalObject, pod *v1.Pod) error {
	for i := range jobname {
		configmap := &v1.ConfigMap{}
		err := ss.Client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: jobname[i]}, configmap)
		if err != nil && !apierrors.IsNotFound(err) {
			klog.V(5).Infof("get configmap error: %v", err)
			return err
		}
		// the configmap doesn't exist
		if apierrors.IsNotFound(err) {
			data := make(map[string]string)
			// get all local objects produced by the required job
			// hostname -> localobject id
			// TODO: if there are lots of localobjects in the same node
			for _, o := range localobjects {
				if (*o).Labels["job"] == jobname[i] {
					data[(*o).Spec.Hostname] = (*o).Spec.ObjectID
				}
			}
			// get all global objects produced by the required job
			// jobname -> globalobject id
			// TODO: if there are lots of globalobjects produced by the same job
			for _, o := range globalobjects {
				if (*o).Labels["job"] == jobname[i] {
					data[jobname[i]] = (*o).Spec.ObjectID
				}
			}
			cm := v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      jobname[i],
					Namespace: namespace,
				},
				Data: data,
			}
			cm.OwnerReferences = pod.GetOwnerReferences()
			if err := ss.Client.Create(ctx, &cm); err != nil {
				klog.V(5).Infof("create configmap error: %v", err)
				return err
			}
			continue
		}
		klog.V(5).Infof("the configmap [%v/%v] exist!", namespace, jobname[i])

	}

	return nil
}

// VineyardScheduling is a plugin that schedules pods that requires vineyard objects as inputs.
type VineyardScheduling struct {
	client.Client
	handle          framework.Handle
	podLister       listerv1.PodLister
	scheduleTimeout *time.Duration
	state           map[string]*SchedulerState
	podRank         map[string]map[string]int64
}

// New initializes a vineyard scheduler
// func New(configuration *runtime.Unknown, handle framework.FrameworkHandle) (framework.Plugin, error) {
func New(client client.Client, config *rest.Config, obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	klog.Info("Initializing the vineyard scheduler plugin ...")
	timeout := Timeout * time.Second
	state := make(map[string]*SchedulerState)
	scheduling := &VineyardScheduling{
		Client:          client,
		handle:          handle,
		podLister:       handle.SharedInformerFactory().Core().V1().Pods().Lister(),
		scheduleTimeout: &timeout,
		state:           state,
		podRank:         map[string]map[string]int64{},
	}
	return scheduling, nil
}

// Name returns name of the plugin. It is used in logs, etc.
func (vs *VineyardScheduling) Name() string {
	return Name
}

// Less compares the priority of two
func (vs *VineyardScheduling) Less(pod1, pod2 *framework.PodInfo) bool {
	prio1 := corev1.PodPriority(pod1.Pod)
	prio2 := corev1.PodPriority(pod2.Pod)
	return prio1 > prio2
}

// Score compute the score for a pod based on the status of required vineyard objects.
func (vs *VineyardScheduling) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	klog.V(5).Infof("scoring for pod %v on node %v", GetNamespacedName(pod), nodeName)

	job, replica, requires, vineyardd, err := vs.GetVineyardInfo(pod)
	if err != nil {
		return 0, framework.NewStatus(framework.Unschedulable, err.Error())
	}

	klog.V(5).Infof("scoring for pod of job %v, with %v replicas (rank %v), and requires %v", job, replica, requires)

	schedulerState := vs.MakeSchedulerStateForNamespace(VineyardSystemNamespace)
	podRank := vs.GetPodRank(pod, replica)
	nodes := vs.GetAllWorkerNodes(vineyardd)

	score, err := schedulerState.Compute(ctx, job, replica, podRank, nodes, requires, nodeName, pod)
	if err != nil {
		return 0, framework.NewStatus(framework.Unschedulable, err.Error())
	}
	if score == 0 {
		return score, framework.NewStatus(framework.Unschedulable, "")
	}
	klog.Infof("score for pod of job %v on node %v is: %v", job, nodeName, score)
	return score, framework.NewStatus(framework.Success, "")
}

// ScoreExtensions of the Score plugin.
func (vs *VineyardScheduling) ScoreExtensions() framework.ScoreExtensions {
	klog.V(5).Infof("ScoreExtensions...")
	return vs
}

// NormalizeScore normalizes the score of all nodes for a pod.
func (vs *VineyardScheduling) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	klog.V(5).Infof("NormalizeScore...")
	// Find highest and lowest scores.
	return framework.NewStatus(framework.Success, "")
}

// Permit only permit runs on the node that has vineyard installed.
func (vs *VineyardScheduling) Permit(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (*framework.Status, time.Duration) {
	klog.V(5).Infof("Permit...")
	return framework.NewStatus(framework.Success, ""), 0
}

// PostBind prints the bind info
func (vs *VineyardScheduling) PostBind(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, nodeName string) {
	klog.V(5).Infof("Bind pod %v on node %v", GetNamespacedName(pod), nodeName)
}

// MakeSchedulerStateForNamespace initializes a state for the given namespace, if not exists.
func (vs *VineyardScheduling) MakeSchedulerStateForNamespace(namespace string) *SchedulerState {
	if _, ok := vs.state[namespace]; !ok {
		state := make(map[string]map[string]string)
		vs.state[namespace] = &SchedulerState{
			Client: vs.Client,
			state:  state,
		}
	}
	return vs.state[namespace]
}

func (vs *VineyardScheduling) getJobName(pod *v1.Pod) (string, error) {
	jobName, exists := pod.Labels[VineyardJobName]
	klog.V(5).Infof("labels: %v", pod.Labels)
	if !exists || jobName == "" {
		return "", fmt.Errorf("Failed to get vineyard job name for %v", GetNamespacedName(pod))
	}
	klog.V(5).Infof("Get job's name: %v", jobName)
	return jobName, nil
}

func (vs *VineyardScheduling) getJobReplica(pod *v1.Pod) (int64, error) {
	klog.V(5).Infof("getJobReplica...")
	// infer from the ownership
	ctx := context.TODO()
	//ctx := context.Background()
	for _, owner := range pod.GetOwnerReferences() {
		name := types.NamespacedName{Namespace: pod.Namespace, Name: owner.Name}
		switch owner.Kind {
		case "ReplicaSet":
			replicaset := &appsv1.ReplicaSet{}
			if err := vs.Get(ctx, name, replicaset); err == nil {
				return int64(*replicaset.Spec.Replicas), nil
			}
		case "DaemonSet":
			daemonset := &appsv1.DaemonSet{}
			if err := vs.Get(ctx, name, daemonset); err == nil {
				return int64(daemonset.Spec.Size()), nil
			}
		case "StatefulSet":
			statefulset := &appsv1.StatefulSet{}
			if err := vs.Get(ctx, name, statefulset); err == nil {
				return int64(*statefulset.Spec.Replicas), nil
			}
		case "Job":
			job := &batchv1.Job{}
			if err := vs.Get(ctx, name, job); err == nil {
				return int64(*job.Spec.Parallelism), nil
			}
		case "CronJob":
			crobjob := &batchv1.CronJob{}
			if err := vs.Get(ctx, name, crobjob); err == nil {
				return int64(crobjob.Spec.Size()), nil
			}
		case "Deployment":
			deployment := &appsv1.Deployment{}
			if err := vs.Get(ctx, name, deployment); err == nil {
				return int64(*deployment.Spec.Replicas), nil
			}
		default:
			klog.V(5).Infof("Unable to infer the job replica, unknown owner kind: %v", owner.Kind)
		}
	}

	return -1, fmt.Errorf("Failed to get vineyard job name for %v", GetNamespacedName(pod))
}

// GetAllWorkerNodes records every worker node which deployed vineyardd.
func (vs *VineyardScheduling) GetAllWorkerNodes(vineyardd string) []string {
	nodes := []string{}

	podList := v1.PodList{}
	option := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{
			"app.kubernetes.io/name":     vineyardd,
			"app.kubernetes.io/instance": "vineyardd",
		}),
	}
	if err := vs.Client.List(context.TODO(), &podList, option); err != nil {
		klog.V(5).Infof("Failed to list all pods with the specific label: %v", err)
	}

	for _, pod := range podList.Items {
		nodes = append(nodes, pod.Spec.NodeName)
	}
	sort.Strings(nodes)
	return nodes
}

// get all required jobs name that separated by '.'
func (vs *VineyardScheduling) getRequiredJob(pod *v1.Pod) ([]string, error) {
	objects, exists := pod.Annotations[VineyardJobRequired]
	if !exists {
		return []string{}, fmt.Errorf("Failed to get the required jobs, please set none if there is no required job")
	}

	klog.V(5).Infof("Get the required jobs: %v", objects)
	if objects == "none" {
		return []string{}, nil
	}
	return strings.Split(objects, "."), nil
}

// GetVineyardInfo requires (job, replica, requires, vineyardd) information of a pod.
func (vs *VineyardScheduling) GetVineyardInfo(pod *v1.Pod) (string, int64, []string, string, error) {
	job, err := vs.getJobName(pod)
	if err != nil {
		return "", 0, nil, "", err
	}
	replica, err := vs.getJobReplica(pod)
	if err != nil {
		return "", 0, nil, "", err
	}
	requires, err := vs.getRequiredJob(pod)
	if err != nil {
		return "", 0, nil, "", err
	}
	vineyardd, exist := pod.Labels[VineyarddName]
	if !exist {
		klog.V(5).Infof("VineyarddName does't exist!")
	}
	return job, replica, requires, vineyardd, nil
}

// GetPodRank returns the rank of this pod
func (vs *VineyardScheduling) GetPodRank(pod *v1.Pod, replica int64) int64 {
	// get workflow's prefix name.
	podName := pod.GetName()
	prefixIndex := strings.LastIndexByte(podName, '-')
	prefixName := podName[:prefixIndex]

	// clean up the pod rank
	if int64(len(vs.podRank[prefixName])) > replica {
		delete(vs.podRank, prefixName)
	}

	rank, prefixExist := vs.podRank[prefixName]
	if !prefixExist {
		m := make(map[string]int64)
		m[podName] = int64(len(rank))
		vs.podRank[prefixName] = m
		return 0
	} else {
		_, podExist := rank[podName]
		if !podExist {
			rank[podName] = int64(len(rank))
		}
	}

	return rank[podName]
}

// GetNamespacedName returns the namespaced name of an kubernetes object.
func GetNamespacedName(object metav1.Object) string {
	return fmt.Sprintf("%v/%v", object.GetNamespace(), object.GetName())
}
