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

	// ingresses, err := clientset.NetworkingV1().Ingresses("").List(context.Background(), v1.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	ingressClasses, err := clientset.NetworkingV1().IngressClasses().List(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, i := range ingressClasses.Items {
		fmt.Printf("%+v\n\n", i.ObjectMeta)
	}

	fmt.Printf("%d\n", len(ingressClasses.Items))
	fmt.Println()

	if len(ingressClasses.Items) > 0 {
		fmt.Printf("%+v\n", ingressClasses.Items[0].ObjectMeta)
	}
}
