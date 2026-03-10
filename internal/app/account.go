package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mindyournow/myn-cli/internal/api"
)

// APIKeyList lists API keys.
func (a *App) APIKeyList(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	keys, err := a.Client.ListAPIKeys(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list API keys: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(keys)
	}
	if len(keys) == 0 {
		return a.Formatter.Println("No API keys.")
	}
	tbl := a.Formatter.NewTable("ID", "NAME", "PREFIX", "ENABLED", "EXPIRES")
	for _, k := range keys {
		enabled := "✓"
		if !k.IsEnabled {
			enabled = "✗"
		}
		tbl.AddRow(k.ID, k.Name, k.Prefix, enabled, k.ExpiresAt)
	}
	tbl.Render()
	return nil
}

// APIKeyCreate creates a new API key.
func (a *App) APIKeyCreate(ctx context.Context, name, description string, scopes []string, expiresAt string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	req := api.CreateAPIKeyRequest{
		Name:        name,
		Description: description,
		Scopes:      scopes,
		ExpiresAt:   expiresAt,
	}
	key, err := a.Client.CreateAPIKey(ctx, req)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to create API key: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(key)
	}
	_ = a.Formatter.Success(fmt.Sprintf("Created API key: %s (%s)", key.Name, key.ID))
	if key.Secret != "" {
		return a.Formatter.Println(fmt.Sprintf("Secret (save this now): %s", key.Secret))
	}
	return nil
}

// APIKeyRevoke deletes an API key.
func (a *App) APIKeyRevoke(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.RevokeAPIKey(ctx, id); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to revoke API key: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Revoked API key %s.", id))
}

// ExportRequest requests a data export.
func (a *App) ExportRequest(ctx context.Context, format string, includes []string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	req := api.ExportRequest{Format: format, Includes: includes}
	export, err := a.Client.RequestExport(ctx, req)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to request export: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(export)
	}
	return a.Formatter.Success(fmt.Sprintf("Export requested: %s (status: %s)", export.ID, export.Status))
}

// ExportList lists data exports.
func (a *App) ExportList(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	exports, err := a.Client.ListExports(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list exports: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(exports)
	}
	if len(exports) == 0 {
		return a.Formatter.Println("No exports.")
	}
	tbl := a.Formatter.NewTable("ID", "FORMAT", "STATUS", "CREATED")
	for _, e := range exports {
		tbl.AddRow(e.ID, e.Format, e.Status, e.CreatedAt)
	}
	tbl.Render()
	return nil
}

// ExportDownload downloads an export to a file.
func (a *App) ExportDownload(ctx context.Context, id, outPath string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	data, err := a.Client.DownloadExport(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to download export: %v", err))
		return err
	}
	if outPath == "" || outPath == "-" {
		_, err = os.Stdout.Write(data)
		return err
	}
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to write file: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Export saved to %s.", outPath))
}

// AccountUsage shows account usage.
func (a *App) AccountUsage(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	usage, err := a.Client.GetAccountUsage(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get usage: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(usage)
	}
	b, err := json.MarshalIndent(usage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}
	return a.Formatter.Println(string(b))
}

// AccountLimits shows subscription limits.
func (a *App) AccountLimits(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	limits, err := a.Client.GetAccountLimits(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get limits: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(limits)
	}
	b, err := json.MarshalIndent(limits, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}
	return a.Formatter.Println(string(b))
}
