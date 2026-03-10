package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/output"
)

// nowDate returns today's date as YYYY-MM-DD.
func nowDate() string {
	return time.Now().Format("2006-01-02")
}

// buildParams creates TaskListParams for a given task type and due-today filter.
func buildParams(taskType string, dueToday bool) api.TaskListParams {
	p := api.TaskListParams{Type: taskType}
	if dueToday {
		p.Today = true
	}
	return p
}

// printTaskList renders a slice of tasks to the formatter's output.
func (a *App) printTaskList(_ context.Context, tasks []api.UnifiedTask, emptyMsg string) error {
	if a.Formatter.JSON {
		type result struct {
			Tasks []api.UnifiedTask `json:"tasks"`
			Count int               `json:"count"`
		}
		return a.Formatter.Print(result{Tasks: tasks, Count: len(tasks)})
	}
	if len(tasks) == 0 {
		return a.Formatter.Println(emptyMsg)
	}
	tbl := a.Formatter.NewTable("", "TITLE", "PRIORITY", "DATE", "TYPE")
	for _, t := range tasks {
		symbol := output.PriorityColored(t.PriorityString(), a.Formatter.NoColor)
		date := formatTaskDate(t.StartDate)
		tbl.AddRow(symbol, t.Title, t.PriorityString(), date, strings.ToLower(t.TaskType))
	}
	tbl.Render()
	return nil
}

// fmtJSON marshals v as indented JSON and prints it.
func (a *App) fmtJSON(v interface{}) error {
	if a.Formatter.JSON {
		return a.Formatter.Print(v)
	}
	// For non-JSON mode, use Println with a basic representation
	return a.Formatter.Println(fmt.Sprintf("%v", v))
}
