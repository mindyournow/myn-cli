package app

// NowList shows current now/focus items.
func (a *App) NowList() error {
	return ErrNotImplemented
}

// NowFocus enters focus mode on a task.
func (a *App) NowFocus(id string) error {
	return ErrNotImplemented
}

// NowComplete completes the current focus task.
func (a *App) NowComplete(id string) error {
	return ErrNotImplemented
}

// NowSnooze snoozes the current focus task.
func (a *App) NowSnooze(duration string) error {
	return ErrNotImplemented
}
