package app

// CompassShow shows today's Compass briefing.
func (a *App) CompassShow() error {
	return ErrNotImplemented
}

// CompassGenerate generates a new Compass briefing.
func (a *App) CompassGenerate() error {
	return ErrNotImplemented
}

// CompassCorrect corrects yesterday's Compass results.
func (a *App) CompassCorrect() error {
	return ErrNotImplemented
}

// CompassComplete marks a Compass item as completed.
func (a *App) CompassComplete(itemID string) error {
	return ErrNotImplemented
}

// CompassStatus shows Compass generation status.
func (a *App) CompassStatus() error {
	return ErrNotImplemented
}

// CompassHistory shows Compass history.
func (a *App) CompassHistory() error {
	return ErrNotImplemented
}
