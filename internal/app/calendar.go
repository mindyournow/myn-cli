package app

// CalendarList lists calendar events.
func (a *App) CalendarList() error {
	return ErrNotImplemented
}

// CalendarAdd adds a calendar event.
func (a *App) CalendarAdd(title string) error {
	return ErrNotImplemented
}

// CalendarDelete deletes a calendar event.
func (a *App) CalendarDelete(id string) error {
	return ErrNotImplemented
}

// CalendarDecline declines a calendar invitation.
func (a *App) CalendarDecline(id string) error {
	return ErrNotImplemented
}

// CalendarSkip skips a recurring calendar event.
func (a *App) CalendarSkip(id string) error {
	return ErrNotImplemented
}
