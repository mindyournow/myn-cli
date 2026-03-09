package app

// TimerStart starts a countdown timer.
func (a *App) TimerStart(duration, name string) error {
	return ErrNotImplemented
}

// TimerAlarm sets an alarm for a specific time.
func (a *App) TimerAlarm(time, name string) error {
	return ErrNotImplemented
}

// TimerPause pauses a timer.
func (a *App) TimerPause(id string) error {
	return ErrNotImplemented
}

// TimerResume resumes a paused timer.
func (a *App) TimerResume(id string) error {
	return ErrNotImplemented
}

// TimerComplete completes a timer early.
func (a *App) TimerComplete(id string) error {
	return ErrNotImplemented
}

// TimerDismiss dismisses an alarm.
func (a *App) TimerDismiss(id string) error {
	return ErrNotImplemented
}

// TimerCount shows active timer count.
func (a *App) TimerCount() error {
	return ErrNotImplemented
}
