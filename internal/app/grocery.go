package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/mindyournow/myn-cli/internal/api"
)

// GroceryList lists grocery items.
func (a *App) GroceryList(ctx context.Context) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	items, err := a.Client.ListGroceryItems(ctx, householdID)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list grocery items: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(items)
	}
	if len(items) == 0 {
		return a.Formatter.Println("Grocery list is empty.")
	}
	tbl := a.Formatter.NewTable("✓", "ITEM", "QTY", "CATEGORY")
	for _, item := range items {
		checked := " "
		if item.IsChecked {
			checked = "✓"
		}
		qty := ""
		if item.Quantity > 0 {
			qty = fmt.Sprintf("%.0f", item.Quantity)
			if item.Unit != "" {
				qty += " " + item.Unit
			}
		}
		tbl.AddRow(checked, item.Name, qty, item.Category)
	}
	tbl.Render()
	return nil
}

// GroceryAdd adds an item to the grocery list.
func (a *App) GroceryAdd(ctx context.Context, name, unit string, quantity float64, category string) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	item, err := a.Client.AddGroceryItem(ctx, householdID, api.AddGroceryItemRequest{
		Name:     name,
		Quantity: quantity,
		Unit:     unit,
		Category: category,
	})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to add grocery item: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(item)
	}
	return a.Formatter.Success(fmt.Sprintf("Added: %s", item.Name))
}

// GroceryAddBulk adds multiple items from a newline/comma-separated list.
func (a *App) GroceryAddBulk(ctx context.Context, itemsStr string) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	var reqs []api.AddGroceryItemRequest
	for _, name := range strings.Split(itemsStr, "\n") {
		name = strings.TrimSpace(name)
		if name != "" {
			reqs = append(reqs, api.AddGroceryItemRequest{Name: name})
		}
	}
	if len(reqs) == 0 {
		return fmt.Errorf("no items provided")
	}
	items, err := a.Client.AddGroceryItemsBulk(ctx, householdID, reqs)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to add grocery items: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(items)
	}
	return a.Formatter.Success(fmt.Sprintf("Added %d items.", len(items)))
}

// GroceryCheck toggles the checked state of a grocery item.
func (a *App) GroceryCheck(ctx context.Context, id string, checked bool) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	item, err := a.Client.CheckGroceryItem(ctx, householdID, id, checked)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to update grocery item: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(item)
	}
	if checked {
		return a.Formatter.Success(fmt.Sprintf("Checked: %s", item.Name))
	}
	return a.Formatter.Success(fmt.Sprintf("Unchecked: %s", item.Name))
}

// GroceryDelete removes a grocery item.
func (a *App) GroceryDelete(ctx context.Context, id string) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	if err := a.Client.DeleteGroceryItem(ctx, householdID, id); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to delete grocery item: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Deleted item %s.", id))
}

// GroceryClear removes all checked items from the list.
func (a *App) GroceryClear(ctx context.Context) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	if err := a.Client.ClearCheckedGroceryItems(ctx, householdID); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to clear grocery list: %v", err))
		return err
	}
	return a.Formatter.Success("Cleared all checked items.")
}

// GroceryConvert converts grocery items to tasks.
func (a *App) GroceryConvert(ctx context.Context, itemIDs []string) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	tasks, err := a.Client.ConvertGroceryToTasks(ctx, householdID, itemIDs)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to convert grocery items: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(tasks)
	}
	return a.Formatter.Success(fmt.Sprintf("Converted %d items to tasks.", len(tasks)))
}
