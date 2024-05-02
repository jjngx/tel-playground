package telemetry

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	_ "github.com/nginxinc/telemetry-exporter/pkg/telemetry"
)

func ClusterID(ctx context.Context, k8sClient client.Reader) (string, error) {
	key := types.NamespacedName{
		Name: meta.NamespaceSystem,
	}
	var kubeNamespace v1.Namespace
	if err := k8sClient.Get(ctx, key, &kubeNamespace); err != nil {
		return "", fmt.Errorf("failed to get kube-system namespace: %w", err)
	}
	return string(kubeNamespace.GetUID()), nil
}

func NodeCount(ctx context.Context, k8sClient client.Reader) (int, error) {
	var nodes v1.NodeList
	if err := k8sClient.List(ctx, &nodes); err != nil {
		return 0, fmt.Errorf("failed to get NodeList: %w", err)
	}
	return len(nodes.Items), nil
}
