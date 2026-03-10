package app

import (
	"context"
	"fmt"

	"github.com/mindyournow/myn-cli/internal/api"
)

// ProjectList lists all projects.
func (a *App) ProjectList(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	projects, err := a.Client.ListProjects(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list projects: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(projects)
	}
	if len(projects) == 0 {
		return a.Formatter.Println("No projects found.")
	}
	tbl := a.Formatter.NewTable("ID", "TITLE", "STATUS")
	for _, p := range projects {
		tbl.AddRow(p.ID, p.Title, p.Status)
	}
	tbl.Render()
	return nil
}

// ProjectShow shows a single project.
func (a *App) ProjectShow(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	project, err := a.Client.GetProject(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get project: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(project)
	}
	_ = a.Formatter.Println(fmt.Sprintf("ID:          %s", project.ID))
	_ = a.Formatter.Println(fmt.Sprintf("Title:       %s", project.Title))
	_ = a.Formatter.Println(fmt.Sprintf("Status:      %s", project.Status))
	if project.Description != "" {
		_ = a.Formatter.Println("")
		return a.Formatter.PrintMarkdown(project.Description)
	}
	return nil
}

// ProjectCreate creates a new project.
func (a *App) ProjectCreate(ctx context.Context, title, description string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	project, err := a.Client.CreateProject(ctx, api.CreateProjectRequest{
		Title:       title,
		Description: description,
	})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to create project: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(project)
	}
	return a.Formatter.Success(fmt.Sprintf("Created project: %s (%s)", project.Title, project.ID))
}
