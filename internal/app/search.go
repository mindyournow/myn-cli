package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/output"
)

// SearchAll performs a unified search.
func (a *App) SearchAll(ctx context.Context, query string, opts SearchOptions) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	results, err := a.Client.Search(ctx, api.SearchParams{
		Query:           query,
		Types:           opts.Types,
		Priorities:      opts.Priorities,
		IncludeArchived: opts.IncludeArchived,
		Limit:           opts.Limit,
		Offset:          opts.Offset,
	})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("search failed: %v", err))
		return err
	}
	if a.Formatter.JSON {
		type searchResult struct {
			Results []api.SearchResult `json:"results"`
			Count   int                `json:"count"`
		}
		return a.Formatter.Print(searchResult{Results: results, Count: len(results)})
	}
	if len(results) == 0 {
		return a.Formatter.Println("No results found.")
	}
	tbl := a.Formatter.NewTable("TYPE", "", "TITLE", "DATE")
	for _, r := range results {
		priority := ""
		if p, ok := r.Priority.(string); ok {
			priority = output.PriorityColored(p, a.Formatter.NoColor)
		}
		tbl.AddRow(strings.ToLower(r.Type), priority, r.Title, r.Date)
	}
	tbl.Render()
	return nil
}

// SearchOptions are filters for search.
type SearchOptions struct {
	Types           []string
	Priorities      []string
	IncludeArchived bool
	Limit           int
	Offset          int
}
