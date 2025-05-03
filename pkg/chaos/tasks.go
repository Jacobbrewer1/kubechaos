package chaos

import (
	"context"
	"crypto/rand"
	"log/slog"
	"math/big"
	"time"

	"github.com/jacobbrewer1/web"
	"github.com/jacobbrewer1/web/logging"
	"k8s.io/client-go/kubernetes"
	listersv1 "k8s.io/client-go/listers/core/v1"
)

// IndefinitePodChaos is a web.AsyncTaskFunc that runs indefinitely, executing pod chaos tasks at regular intervals.
func IndefinitePodChaos(
	l *slog.Logger,
	kubeClient kubernetes.Interface,
	podLister listersv1.PodLister,
) web.AsyncTaskFunc {
	return func(ctx context.Context) {
		l.Info("starting indefinite task")

		for {
			randInt, err := rand.Int(rand.Reader, big.NewInt(10000))
			if err != nil {
				l.Error("error generating random number", slog.String(logging.KeyError, err.Error()))
				return
			}

			select {
			case <-ctx.Done():
				l.Debug("context cancelled, stopping indefinite task")
				return
			case <-time.After(time.Duration(randInt.Int64()) * time.Millisecond):
				if err := executePodChaos(ctx, l, kubeClient, podLister); err != nil {
					l.Error("error executing pod chaos", slog.String(logging.KeyError, err.Error()))
				}
			}
		}
	}
}
