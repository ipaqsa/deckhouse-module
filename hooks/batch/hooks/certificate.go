package hooks

import (
	"context"
	"log/slog"

	"github.com/deckhouse/module-sdk/pkg"
	objectpatch "github.com/deckhouse/module-sdk/pkg/object-patch"
	"github.com/deckhouse/module-sdk/pkg/registry"
)

var _ = registry.RegisterFunc(tlsSecretHookConfig, handlerHook)

var tlsSecretHookConfig = &pkg.HookConfig{
	OnBeforeHelm: &pkg.OrderedConfig{Order: 10},
	Kubernetes: []pkg.KubernetesConfig{
		{
			Name:       "tlsSecret",
			APIVersion: "v1",
			Kind:       "Secret",
			JqFilter:   `{"name": .metadata.name, "crt": .data."ca.crt"}`,

			NameSelector: &pkg.NameSelector{
				MatchNames: []string{"webhook-handler-certs"},
			},

			NamespaceSelector: &pkg.NamespaceSelector{
				NameSelector: &pkg.NameSelector{
					MatchNames: []string{"d8-system"},
				},
			},
		},
	},

	Queue: "/xmodule",
}

type cert struct {
	Name     string `json:"name"`
	CertData string `json:"crt"`
}

func handlerHook(_ context.Context, input *pkg.HookInput) error {
	data, err := objectpatch.UnmarshalToStruct[cert](input.Snapshots, "tlsSecret")
	if err != nil {
		return err
	}

	input.Logger.Info("secret hook done", slog.Any("data", data))

	return nil
}
