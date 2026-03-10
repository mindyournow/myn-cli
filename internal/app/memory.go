package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/mindyournow/myn-cli/internal/api"
)

// MemoryList lists all memories.
func (a *App) MemoryList(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	memories, err := a.Client.ListMemories(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list memories: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(memories)
	}
	if len(memories) == 0 {
		return a.Formatter.Println("No memories found.")
	}
	tbl := a.Formatter.NewTable("ID", "CONTENT", "TAGS")
	for _, m := range memories {
		tags := strings.Join(m.Tags, ", ")
		content := m.Content
		if len(content) > 60 {
			content = content[:57] + "..."
		}
		tbl.AddRow(m.ID, content, tags)
	}
	tbl.Render()
	return nil
}

// MemoryShow displays a single memory.
func (a *App) MemoryShow(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	memory, err := a.Client.GetMemory(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get memory: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(memory)
	}
	_ = a.Formatter.Println(fmt.Sprintf("ID:      %s", memory.ID))
	_ = a.Formatter.Println(fmt.Sprintf("Tags:    %s", strings.Join(memory.Tags, ", ")))
	_ = a.Formatter.Println("")
	return a.Formatter.PrintMarkdown(memory.Content)
}

// MemoryAdd adds a new memory.
func (a *App) MemoryAdd(ctx context.Context, content string, tags []string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	memory, err := a.Client.AddMemory(ctx, api.CreateMemoryRequest{
		Content: content,
		Tags:    tags,
	})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to add memory: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(memory)
	}
	return a.Formatter.Success(fmt.Sprintf("Memory added: %s", memory.ID))
}

// MemoryUpdate updates a memory.
func (a *App) MemoryUpdate(ctx context.Context, id, content string, tags []string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	memory, err := a.Client.UpdateMemory(ctx, id, api.CreateMemoryRequest{
		Content: content,
		Tags:    tags,
	})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to update memory: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(memory)
	}
	return a.Formatter.Success(fmt.Sprintf("Memory updated: %s", memory.ID))
}

// MemorySearch searches memories.
func (a *App) MemorySearch(ctx context.Context, query string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	memories, err := a.Client.SearchMemories(ctx, query)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("search failed: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(memories)
	}
	if len(memories) == 0 {
		return a.Formatter.Println("No memories matched.")
	}
	tbl := a.Formatter.NewTable("ID", "CONTENT", "TAGS")
	for _, m := range memories {
		content := m.Content
		if len(content) > 60 {
			content = content[:57] + "..."
		}
		tbl.AddRow(m.ID, content, strings.Join(m.Tags, ", "))
	}
	tbl.Render()
	return nil
}

// MemoryDelete deletes a memory.
func (a *App) MemoryDelete(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.DeleteMemory(ctx, id); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to delete memory: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Memory deleted: %s", id))
}

// MemoryExport exports all memories.
func (a *App) MemoryExport(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	data, err := a.Client.ExportMemories(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to export memories: %v", err))
		return err
	}
	_, writeErr := fmt.Print(string(data))
	return writeErr
}
