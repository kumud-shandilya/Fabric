package main

import (
	"github.com/getporter/FabricNew/pkg/FabricNew"
	"github.com/spf13/cobra"
)

func buildUpgradeCommand(m *FabricNew.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Execute the invoke functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute(cmd.Context())
		},
	}
	return cmd
}
