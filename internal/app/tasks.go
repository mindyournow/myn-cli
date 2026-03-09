package app

import (
	"fmt"
)

// TaskList lists all tasks.
func (a *App) TaskList() error {
	fmt.Println("Task list not yet implemented.")
	return nil
}

// TaskAdd adds a new task.
func (a *App) TaskAdd(title string) error {
	fmt.Printf("Task add not yet implemented: %s\n", title)
	return nil
}

// TaskShow shows task details.
func (a *App) TaskShow(id string) error {
	fmt.Printf("Task show not yet implemented: %s\n", id)
	return nil
}

// TaskEdit edits a task.
func (a *App) TaskEdit(id string) error {
	fmt.Printf("Task edit not yet implemented: %s\n", id)
	return nil
}

// TaskDelete deletes a task.
func (a *App) TaskDelete(id string) error {
	fmt.Printf("Task delete not yet implemented: %s\n", id)
	return nil
}

// TaskDone marks a task as done.
func (a *App) TaskDone(id string) error {
	fmt.Printf("Task done not yet implemented: %s\n", id)
	return nil
}

// TaskRestore restores a deleted task.
func (a *App) TaskRestore(id string) error {
	fmt.Printf("Task restore not yet implemented: %s\n", id)
	return nil
}

// TaskArchive archives a task.
func (a *App) TaskArchive(id string) error {
	fmt.Printf("Task archive not yet implemented: %s\n", id)
	return nil
}

// TaskSnooze snoozes a task.
func (a *App) TaskSnooze(id, duration string) error {
	fmt.Printf("Task snooze not yet implemented: %s for %s\n", id, duration)
	return nil
}

// TaskMove moves a task to a different project.
func (a *App) TaskMove(id, projectID string) error {
	fmt.Printf("Task move not yet implemented: %s to %s\n", id, projectID)
	return nil
}

// TaskBatch performs batch operations on tasks.
func (a *App) TaskBatch(operation string, ids []string) error {
	fmt.Printf("Task batch not yet implemented: %s on %v\n", operation, ids)
	return nil
}
