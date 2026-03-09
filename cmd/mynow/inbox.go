package main

import (
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/spf13/cobra"
)

func newInboxCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inbox",
		Short: "Manage inbox items",
	}

	cmd.AddCommand(
		newInboxListCmd(a),
		newInboxAddCmd(a),
		newInboxProcessCmd(a),
		newInboxClearCmd(a),
	)

	return cmd
}

func newInboxListCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List inbox items",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.InboxList()
		},
	}
}

func newInboxAddCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "add <title>",
		Short: "Add an item to the inbox",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.InboxAdd(args[0])
		},
	}
}

func newInboxProcessCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "process [id]",
		Short: "Process inbox items (interactive if no ID)",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			}
			return a.InboxProcess(id)
		},
	}
}

func newInboxClearCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Clear all inbox items",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.InboxClear()
		},
	}
}
