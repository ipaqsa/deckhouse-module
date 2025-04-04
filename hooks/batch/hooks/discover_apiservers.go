package hooks

import (
	"context"
	"log/slog"

	"github.com/deckhouse/module-sdk/pkg"
	objectpatch "github.com/deckhouse/module-sdk/pkg/object-patch"
	"github.com/deckhouse/module-sdk/pkg/registry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = registry.RegisterFunc(config, handler)

var config = &pkg.HookConfig{
	Kubernetes: []pkg.KubernetesConfig{
		{
			Name:       "apiservers",
			APIVersion: "v1",
			Kind:       "Pod",
			NamespaceSelector: &pkg.NamespaceSelector{
				NameSelector: &pkg.NameSelector{
					MatchNames: []string{"kube-system"},
				},
			},
			LabelSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"component": "kube-apiserver"},
			},
			JqFilter: ".metadata.name",
		},
	},
}

func handler(_ context.Context, input *pkg.HookInput) error {
	podNames, err := objectpatch.UnmarshalToStruct[string](input.Snapshots, "apiservers")
	if err != nil {
		return err
	}

	input.Logger.Info("discover api servers hook done", slog.Any("podNames", podNames))

	input.Values.Set("xmodule.internal.apiServers", podNames)

	return nil
}
