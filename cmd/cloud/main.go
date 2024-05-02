package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, n := range nodes.Items {
		fmt.Printf("Node Name   => %+v\n", n.Name)
		fmt.Printf("Provider ID => %+v\n", n.Spec.ProviderID)
	}
}

// fmt.Printf("%+v\n", r)

// fmt.Println("========= deployments =========")

// dpls, err := clientset.AppsV1().Deployments("nginx-ingress").List(context.Background(), metav1.ListOptions{})
// if err != nil {
// 	panic(err)
// }

// for _, d := range dpls.Items {
// 	fmt.Printf("%+v\n", d.Kind)
// }

// fmt.Println("=========")

// // for _, d := range dpls.Items {
// // 	fmt.Printf("%+v\n", d.ObjectMeta)
// // }

// // // for _, i := range rs.Items {
// // // 	fmt.Printf("%+v\n", i.OwnerReferences[0].Name)
// // // }

// // fmt.Println("========= ABC =========")

// // // 1. Get Pod based on name and namespace
// // // 2. Get kind from pod owner ref that is "ReplicSet"
// // // 3. Get replica set

// // pod, err := clientset.CoreV1().Pods("nginx-ingress").Get(context.Background(), "nginx-ingress-9b9b8f76b-pv4kv", metav1.GetOptions{})
// // if err != nil {
// // 	panic(err)
// // }

// fmt.Println("========= Pod =========")

// p, err := clientset.CoreV1().Pods("nginx-ingress").Get(context.Background(), "nginx-ingress-9b9b8f76b-pv4kv", metav1.GetOptions{})
// if err != nil {
// 	panic(err)
// }
// fmt.Printf("%+v\n", p.OwnerReferences[0].Kind)
// fmt.Printf("%+v\n", p.OwnerReferences[0].Name)

// fmt.Println("========= ReplicaSet =========")
// rs, err := clientset.AppsV1().ReplicaSets("nginx-ingress").Get(context.Background(), p.OwnerReferences[0].Name, metav1.GetOptions{})
// if err != nil {
// 	panic(err)
// }

// fmt.Printf("%+v\n", *rs.Spec.Replicas)

// func getPodReplicaSet(
// 	ctx context.Context,
// 	k8sClient client.Reader,
// 	podNSName types.NamespacedName,
// ) (*appsv1.ReplicaSet, error) {
// 	var pod v1.Pod
// 	if err := k8sClient.Get(
// 		ctx,
// 		types.NamespacedName{Namespace: podNSName.Namespace, Name: podNSName.Name},
// 		&pod,
// 	); err != nil {
// 		return nil, fmt.Errorf("failed to get NGF Pod: %w", err)
// 	}

// 	podOwnerRefs := pod.GetOwnerReferences()
// 	if len(podOwnerRefs) != 1 {
// 		return nil, fmt.Errorf("expected one owner reference of the NGF Pod, got %d", len(podOwnerRefs))
// 	}

// 	if podOwnerRefs[0].Kind != "ReplicaSet" {
// 		return nil, fmt.Errorf("expected pod owner reference to be ReplicaSet, got %s", podOwnerRefs[0].Kind)
// 	}

// 	var replicaSet appsv1.ReplicaSet
// 	if err := k8sClient.Get(
// 		ctx,
// 		types.NamespacedName{Namespace: podNSName.Namespace, Name: podOwnerRefs[0].Name},
// 		&replicaSet,
// 	); err != nil {
// 		return nil, fmt.Errorf("failed to get NGF Pod's ReplicaSet: %w", err)
// 	}

// 	return &replicaSet, nil
// }

// func getReplicas(replicaSet *appsv1.ReplicaSet) (int, error) {
// 	if replicaSet.Spec.Replicas == nil {
// 		return 0, errors.New("replica set replicas was nil")
// 	}

// 	return int(*replicaSet.Spec.Replicas), nil
// }
