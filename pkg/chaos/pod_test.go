package chaos

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
)

func testablePod(t *testing.T) *corev1.Pod {
	t.Helper()
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-namespace",
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}
}

func Test_KillPod(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		pod := testablePod(t)
		kubeClient := fake.NewClientset(pod)

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		err := killPod(ctx, kubeClient, pod.Name, pod.Namespace)
		require.NoError(t, err)
	})

	t.Run("pod already terminating", func(t *testing.T) {
		t.Parallel()
		pod := testablePod(t)
		pod.DeletionTimestamp = &metav1.Time{}
		kubeClient := fake.NewClientset(pod)

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		err := killPod(ctx, kubeClient, pod.Name, pod.Namespace)
		require.Error(t, err)
		require.EqualError(t, err, "pod is already terminating")
	})

	t.Run("pod not found", func(t *testing.T) {
		t.Parallel()
		kubeClient := fake.NewClientset()

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		err := killPod(ctx, kubeClient, "non-existent-pod", "test-namespace")
		require.Error(t, err)
		require.Contains(t, err.Error(), "could not get pod")
	})
}

func Test_GetRandomPod(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		pod := testablePod(t)

		kubeClient := fake.NewClientset(pod)
		informerFactory := informers.NewSharedInformerFactory(kubeClient, 5*time.Millisecond)
		podLister := informerFactory.Core().V1().Pods().Lister()
		informerFactory.Start(ctx.Done())
		informerFactory.WaitForCacheSync(ctx.Done())

		randomPod, err := getRandomPod(podLister)
		require.NoError(t, err)
		require.Equal(t, pod.Name, randomPod.Name)
	})

	t.Run("success - multiple pods", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		pod1 := testablePod(t)
		pod2 := testablePod(t)
		pod2.Name = "test-pod-2"

		kubeClient := fake.NewClientset(pod1, pod2)
		informerFactory := informers.NewSharedInformerFactory(kubeClient, 5*time.Millisecond)
		podLister := informerFactory.Core().V1().Pods().Lister()
		informerFactory.Start(ctx.Done())
		informerFactory.WaitForCacheSync(ctx.Done())

		randomPod, err := getRandomPod(podLister)
		require.NoError(t, err)
		require.Contains(t, []string{pod1.Name, pod2.Name}, randomPod.Name)
	})

	t.Run("no pods found", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		kubeClient := fake.NewClientset()
		informerFactory := informers.NewSharedInformerFactory(kubeClient, 5*time.Millisecond)
		podLister := informerFactory.Core().V1().Pods().Lister()
		informerFactory.Start(ctx.Done())
		informerFactory.WaitForCacheSync(ctx.Done())

		randomPod, err := getRandomPod(podLister)
		require.Error(t, err)
		require.Nil(t, randomPod)
	})
}
