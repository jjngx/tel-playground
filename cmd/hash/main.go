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

	fmt.Println("========= deployments =========")

	dpls, err := clientset.AppsV1().Deployments("nginx-ingress").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, d := range dpls.Items {
		fmt.Printf("%+v\n", d.Kind)
	}

	fmt.Println("========= Pod =========")

	p, err := clientset.CoreV1().Pods("nginx-ingress").Get(context.Background(), "nginx-ingress-9b9b8f76b-gcrqr", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", p.OwnerReferences[0].Kind)
	fmt.Printf("%+v\n", p.OwnerReferences[0].Name)

	fmt.Println("========= ReplicaSet =========")
	rs, err := clientset.AppsV1().ReplicaSets("nginx-ingress").Get(context.Background(), p.OwnerReferences[0].Name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", rs)

	fmt.Printf("%+v\n", *rs.Spec.Replicas)

	fmt.Println("========= Deployments =========")
	//rs, err := clientset.AppsV1().ReplicaSets("nginx-ingress").Get(context.Background(), p.OwnerReferences[0].Name, metav1.GetOptions{})
	dp, err := clientset.AppsV1().Deployments("nginx-ingress").Get(context.Background(), "nginx-ingress", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", dp.UID)
	fmt.Printf("%+v\n", dp.Name)

	fmt.Printf("%+v\n", dp.Namespace)

}
