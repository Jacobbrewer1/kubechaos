package chaos

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_ExecutePodChaos(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		l := slog.New(slog.DiscardHandler)

		pod := testablePod(t)
		kubeClient := fake.NewClientset(pod)
		informerFactory := informers.NewSharedInformerFactory(kubeClient, 10*time.Millisecond)
		podLister := informerFactory.Core().V1().Pods().Lister()
		informerFactory.Start(ctx.Done())
		informerFactory.WaitForCacheSync(ctx.Done())

		err := executePodChaos(ctx, l, kubeClient, podLister)
		require.NoError(t, err)
		require.Eventually(t, func() bool {
			pod, err := kubeClient.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
			if err != nil {
				return true
			}
			return pod.DeletionTimestamp != nil
		}, 5*time.Second, 100*time.Millisecond, "pod was not deleted in time")
	})

	t.Run("no pods found", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		l := slog.New(slog.DiscardHandler)

		kubeClient := fake.NewClientset()
		informerFactory := informers.NewSharedInformerFactory(kubeClient, 10*time.Millisecond)
		podLister := informerFactory.Core().V1().Pods().Lister()
		informerFactory.Start(ctx.Done())
		informerFactory.WaitForCacheSync(ctx.Done())

		err := executePodChaos(ctx, l, kubeClient, podLister)
		require.Error(t, err)
		require.EqualError(t, err, "could not get random pod: no pods found")
	})
}
