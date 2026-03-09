package app

// GroceryList lists grocery items.
func (a *App) GroceryList() error {
	return ErrNotImplemented
}

// GroceryAdd adds an item to the grocery list.
func (a *App) GroceryAdd(item string) error {
	return ErrNotImplemented
}

// GroceryAddBulk adds multiple items at once.
func (a *App) GroceryAddBulk(items []string) error {
	return ErrNotImplemented
}

// GroceryCheck checks off a grocery item.
func (a *App) GroceryCheck(itemID string) error {
	return ErrNotImplemented
}

// GroceryDelete deletes a grocery item.
func (a *App) GroceryDelete(itemID string) error {
	return ErrNotImplemented
}

// GroceryClear clears all checked items.
func (a *App) GroceryClear() error {
	return ErrNotImplemented
}

// GroceryConvert converts a grocery item to a task.
func (a *App) GroceryConvert(itemID string) error {
	return ErrNotImplemented
}
