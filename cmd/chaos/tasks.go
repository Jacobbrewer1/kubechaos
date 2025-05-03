package main

import (
	"context"

	"github.com/jacobbrewer1/kubechaos/pkg/chaos"
	"github.com/jacobbrewer1/web/logging"
)

func (a *App) podChaos(ctx context.Context) {
	podChaosTask := chaos.IndefinitePodChaos(
		logging.LoggerWithComponent(a.base.Logger(), "pod-chaos"),
		a.base.KubeClient(),
		a.base.PodLister(),
	)
	podChaosTask(ctx)
}
