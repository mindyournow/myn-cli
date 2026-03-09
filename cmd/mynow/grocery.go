package main

import (
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/spf13/cobra"
)

func newGroceryCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grocery",
		Short: "Manage grocery lists",
		Aliases: []string{"groceries", "list", "shop"},
	}

	cmd.AddCommand(
		newGroceryListCmd(a),
		newGroceryAddCmd(a),
		newGroceryAddBulkCmd(a),
		newGroceryCheckCmd(a),
		newGroceryDeleteCmd(a),
		newGroceryClearCmd(a),
		newGroceryConvertCmd(a),
	)

	return cmd
}

func newGroceryListCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List grocery items",
		Aliases: []string{"ls", "show"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.GroceryList()
		},
	}
}

func newGroceryAddCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <item>",
		Short: "Add an item to the grocery list",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.GroceryAdd(args[0])
		},
	}
	cmd.Flags().String("category", "", "Item category")
	return cmd
}

func newGroceryAddBulkCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "add-bulk <items...>",
		Short: "Add multiple items at once",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.GroceryAddBulk(args)
		},
	}
}

func newGroceryCheckCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "check <item-id>",
		Short: "Check off a grocery item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.GroceryCheck(args[0])
		},
	}
}

func newGroceryDeleteCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <item-id>",
		Short: "Delete a grocery item",
		Aliases: []string{"rm", "remove"},
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.GroceryDelete(args[0])
		},
	}
}

func newGroceryClearCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Clear all checked items",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.GroceryClear()
		},
	}
}

func newGroceryConvertCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "convert <item-id>",
		Short: "Convert a grocery item to a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.GroceryConvert(args[0])
		},
	}
}
