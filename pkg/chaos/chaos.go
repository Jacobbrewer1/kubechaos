package chaos

import (
	"context"
	"fmt"
	"log/slog"

	"k8s.io/client-go/kubernetes"
	listersv1 "k8s.io/client-go/listers/core/v1"

	logkeys "github.com/jacobbrewer1/kubechaos/pkg/logging"
	"github.com/jacobbrewer1/web/logging"
)

// executePodChaos executes a pod chaos task that randomly kills a pod in the cluster.
func executePodChaos(
	ctx context.Context,
	l *slog.Logger,
	kubeClient kubernetes.Interface,
	podLister listersv1.PodLister,
) error {
	randomPod, err := getRandomPod(podLister)
	if err != nil {
		return fmt.Errorf("could not get random pod: %w", err)
	}

	l.Debug("random pod selected",
		slog.String(logging.KeyName, randomPod.Name),
		slog.String(logkeys.KeyNamespace, randomPod.Namespace),
	)
	if err := killPod(ctx, kubeClient, randomPod.Name, randomPod.Namespace); err != nil {
		return fmt.Errorf("could not kill pod: %w", err)
	}
	l.Info("pod killed",
		slog.String(logging.KeyName, randomPod.Name),
		slog.String(logkeys.KeyNamespace, randomPod.Namespace),
	)
	return nil
}
