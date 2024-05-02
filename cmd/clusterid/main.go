package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	ns, err := clientset.CoreV1().Namespaces().Get(context.Background(), "kube-system", v1.GetOptions{})
	//cluster, err := clientset.CoreV1().Services("kube-system").Get(context.Background(), "kube-dns", v1.GetOptions{})
	//cluster, err := clientset.CoreV1().Services("kube-system").List(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", ns)

	//telemetry.ClusterID(context.Background(), client.Client)

	fmt.Println("========= ClusterID =========")
	cluster, err := clientset.CoreV1().Namespaces().Get(context.Background(), "kube-system", v1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", cluster.UID)
}
