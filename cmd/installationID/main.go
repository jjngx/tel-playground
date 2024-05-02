package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

	fmt.Println("========= DEPLOYMENT =========")
	// test for deployment
	// inst, err := InstallationIDForDeployment(context.Background(), clientset)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(inst)

	// test for daemonset

	fmt.Println("========= DAEMONSET =========")
	inst, err := InstallationIDForDaemonSet(context.Background(), clientset)
	if err != nil {
		panic(err)
	}

	fmt.Println(inst)
}

func InstallationIDForDeployment(ctx context.Context, cnf *kubernetes.Clientset) (string, error) {
	// Env Vars come from the container
	podNS := "nginx-ingress"
	podName := "nginx-ingress-9b9b8f76b-hzmkv"

	// ========= ACTIONS =========
	// 1) get pod name and owner references => replicaset
	// 2) get replicaset name and owner references => deployment
	// 3) get deployment name
	// ===========================

	// Get the pod info
	fmt.Println("========= Pod =========")

	pod, err := cnf.CoreV1().Pods(podNS).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	podRef := pod.GetOwnerReferences()
	if len(podRef) != 1 {
		return "", fmt.Errorf("expected pod owner reference to be 1, got %d", len(podRef))
	}
	fmt.Printf("kind => %+v\n", podRef[0].Kind)
	fmt.Printf("name => %+v\n", podRef[0].Name)

	rs, err := cnf.AppsV1().ReplicaSets(podNS).Get(ctx, podRef[0].Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	rsRef := rs.GetOwnerReferences()
	if len(rsRef) != 1 {
		return "", fmt.Errorf("expected replicaset owner reference to be 1, got %d", len(rsRef))
	}

	fmt.Printf("RS ref kind => %+v\n", rsRef[0].Kind)
	fmt.Printf("RS ref name => %+v\n", rsRef[0].Name)

	deploymentName := rsRef[0].Name

	// var instID string
	// switch podRef[0].Kind {
	// case "ReplicaSet":
	// 	replicaName := "replica123"
	// 	instID = fmt.Sprintf("%s/%s/%s/deployment", "CID", podNS, replicaName)
	// case "DaemonSet":
	// 	daemonSetName := "daemonSet123"
	// 	instID = fmt.Sprintf("%s/%s/%s/daemonset", "CID", podNS, daemonSetName)
	// default:
	// 	return "", fmt.Errorf("cannot generate InstallationID, expected pod owner reference to be ReplicaSet or DeamonSet, got %s", podRef[0].Kind)
	// }

	instID := fmt.Sprintf("%s/%s/%s/deployment", "CID", podNS, deploymentName)

	return makeHash(instID)

	// h := md5.New()
	// _, err = h.Write([]byte(instID))
	// if err != nil {
	// 	return "", err
	// }
	// return hex.EncodeToString(h.Sum(nil)), nil
}

func InstallationIDForDaemonSet(ctx context.Context, cnf *kubernetes.Clientset) (string, error) {
	// Env Vars come from the container
	podNS := "nginx-ingress"
	podName := "nginx-ingress-test-x6w95"

	// ========= ACTIONS =========
	// 1) get pod name and owner references => replicaset
	// 2) get replicaset name and owner references => deployment
	// 3) get deployment name
	// ===========================

	// Get the pod info
	fmt.Println("========= Pod =========")

	pod, err := cnf.CoreV1().Pods(podNS).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	podRef := pod.GetOwnerReferences()
	if len(podRef) != 1 {
		return "", fmt.Errorf("expected pod owner reference to be 1, got %d", len(podRef))
	}
	fmt.Printf("kind => %+v\n", podRef[0].Kind)
	fmt.Printf("name => %+v\n", podRef[0].Name)

	// ds, err := cnf.AppsV1().DaemonSets(podNS).Get(ctx, podRef[0].Name, metav1.GetOptions{})
	// if err != nil {
	// 	return "", err
	// }

	// dsRef := ds.GetOwnerReferences()
	// if len(dsRef) != 1 {
	// 	return "", fmt.Errorf("expected daemonset owner reference to be 1, got %d", len(dsRef))
	// }

	fmt.Printf("DS ref kind => %+v\n", podRef[0].Kind)
	fmt.Printf("DS ref name => %+v\n", podRef[0].Name)

	daemonSetName := podRef[0].Name

	// var instID string
	// switch podRef[0].Kind {
	// case "ReplicaSet":
	// 	replicaName := "replica123"
	// 	instID = fmt.Sprintf("%s/%s/%s/deployment", "CID", podNS, replicaName)
	// case "DaemonSet":
	// 	daemonSetName := "daemonSet123"
	// 	instID = fmt.Sprintf("%s/%s/%s/daemonset", "CID", podNS, daemonSetName)
	// default:
	// 	return "", fmt.Errorf("cannot generate InstallationID, expected pod owner reference to be ReplicaSet or DeamonSet, got %s", podRef[0].Kind)
	// }

	instID := fmt.Sprintf("%s/%s/%s/daemonset", "CID", podNS, daemonSetName)
	return instID, nil

	// h := md5.New()
	// _, err = h.Write([]byte(instID))
	// if err != nil {
	// 	return "", err
	// }
	// return hex.EncodeToString(h.Sum(nil)), nil
}

func makeHash(identifier string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(identifier))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
