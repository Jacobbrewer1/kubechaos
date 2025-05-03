package chaos

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	listersv1 "k8s.io/client-go/listers/core/v1"
)

// getRandomPod selects a random pod from the given pod lister.
func getRandomPod(podLister listersv1.PodLister) (*corev1.Pod, error) {
	pods, err := podLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("could not list pods: %w", err)
	}

	if len(pods) == 0 {
		return nil, errors.New("no pods found")
	}

	// Select a random pod from the list
	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(pods))))
	if err != nil {
		return nil, fmt.Errorf("could not generate random index: %w", err)
	}
	return pods[randomIndex.Int64()], nil
}

// killPod deletes a pod in the given namespace.
func killPod(
	ctx context.Context,
	kubeClient kubernetes.Interface,
	podName, namespace string,
) error {
	pod, err := kubeClient.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("could not get pod: %w", err)
	}

	if pod.DeletionTimestamp != nil {
		return errors.New("pod is already terminating")
	}

	if err := kubeClient.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("could not delete pod: %w", err)
	}
	return nil
}
