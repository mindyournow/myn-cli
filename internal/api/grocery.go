package api

import (
	"context"
	"fmt"
)

// GroceryItem represents a grocery list item.
type GroceryItem struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Quantity  float64 `json:"quantity,omitempty"`
	Unit      string  `json:"unit,omitempty"`
	IsChecked bool    `json:"isChecked"`
	Category  string  `json:"category,omitempty"`
}

// ListGroceryItems fetches the grocery list for a household.
func (c *Client) ListGroceryItems(ctx context.Context, householdID string) ([]GroceryItem, error) {
	resp, err := c.Get(ctx, "/api/v1/households/"+householdID+"/grocery-list", nil)
	if err != nil {
		return nil, err
	}
	var items []GroceryItem
	if err := resp.DecodeJSON(&items); err != nil {
		return nil, fmt.Errorf("failed to parse grocery list: %w", err)
	}
	return items, nil
}

// AddGroceryItemRequest is the body for adding a grocery item.
type AddGroceryItemRequest struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity,omitempty"`
	Unit     string  `json:"unit,omitempty"`
	Category string  `json:"category,omitempty"`
}

// AddGroceryItem adds a single item to the grocery list.
func (c *Client) AddGroceryItem(ctx context.Context, householdID string, req AddGroceryItemRequest) (*GroceryItem, error) {
	resp, err := c.Post(ctx, "/api/v1/households/"+householdID+"/grocery-list/items", req)
	if err != nil {
		return nil, err
	}
	var item GroceryItem
	if err := resp.DecodeJSON(&item); err != nil {
		return nil, fmt.Errorf("failed to parse grocery item: %w", err)
	}
	return &item, nil
}

// AddGroceryItemsBulk adds multiple items to the grocery list at once.
func (c *Client) AddGroceryItemsBulk(ctx context.Context, householdID string, items []AddGroceryItemRequest) ([]GroceryItem, error) {
	resp, err := c.Post(ctx, "/api/v1/households/"+householdID+"/grocery-list/items/bulk", items)
	if err != nil {
		return nil, err
	}
	var result []GroceryItem
	if err := resp.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("failed to parse bulk grocery items: %w", err)
	}
	return result, nil
}

// CheckGroceryItem toggles the checked state of a grocery item.
func (c *Client) CheckGroceryItem(ctx context.Context, householdID, itemID string, checked bool) (*GroceryItem, error) {
	resp, err := c.Patch(ctx, "/api/v1/households/"+householdID+"/grocery-list/items/"+itemID,
		map[string]bool{"isChecked": checked})
	if err != nil {
		return nil, err
	}
	var item GroceryItem
	if err := resp.DecodeJSON(&item); err != nil {
		return nil, fmt.Errorf("failed to parse grocery item: %w", err)
	}
	return &item, nil
}

// DeleteGroceryItem removes a grocery item.
func (c *Client) DeleteGroceryItem(ctx context.Context, householdID, itemID string) error {
	_, err := c.Delete(ctx, "/api/v1/households/"+householdID+"/grocery-list/"+itemID)
	return err
}

// ClearCheckedGroceryItems removes all checked items from the grocery list.
func (c *Client) ClearCheckedGroceryItems(ctx context.Context, householdID string) error {
	_, err := c.Delete(ctx, "/api/v1/households/"+householdID+"/grocery-list/checked")
	return err
}

// ConvertGroceryToTasks converts grocery items to unified tasks.
func (c *Client) ConvertGroceryToTasks(ctx context.Context, householdID string, itemIDs []string) ([]UnifiedTask, error) {
	resp, err := c.Post(ctx, "/api/v1/households/"+householdID+"/grocery-list/convert-to-tasks",
		map[string][]string{"ids": itemIDs})
	if err != nil {
		return nil, err
	}
	var tasks []UnifiedTask
	if err := resp.DecodeJSON(&tasks); err != nil {
		return nil, fmt.Errorf("failed to parse converted tasks: %w", err)
	}
	return tasks, nil
}
