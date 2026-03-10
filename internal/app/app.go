package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/auth"
	"github.com/mindyournow/myn-cli/internal/config"
	"github.com/mindyournow/myn-cli/internal/output"
)

// App is the central application struct shared by CLI and TUI.
type App struct {
	Config     *config.Config
	Client     *api.Client
	Keyring    *auth.Keyring
	KeyStore   *auth.KeyStore
	TokenCache *auth.TokenCache
	Formatter  *output.Formatter
}

// New creates a new App instance using environment variables and defaults.
// Returns an error if configuration cannot be loaded.
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return NewWithConfig(cfg), nil
}

// NewWithConfig creates an App from an already-loaded Config.
// Used when flags (e.g., --api-url) override the loaded configuration (BUG-3 fix).
func NewWithConfig(cfg *config.Config) *App {
	fileKeyring := auth.NewKeyring(cfg.ConfigDir)
	keyStore := auth.NewKeyStore(fileKeyring, cfg.Auth.Keyring)
	oauthClient := auth.NewOAuthClient(cfg.BaseURL, keyStore)
	tokenCache := auth.NewTokenCache(keyStore, oauthClient)

	return &App{
		Config:     cfg,
		Client:     api.NewClient(cfg.BaseURL),
		Keyring:    fileKeyring,
		KeyStore:   keyStore,
		TokenCache: tokenCache,
		Formatter:  output.NewFormatter(false, false, false),
	}
}

// SetFormatter sets the output formatter for the app.
func (a *App) SetFormatter(f *output.Formatter) {
	a.Formatter = f
}

// Login performs OAuth PKCE authentication.
func (a *App) Login(ctx context.Context, device bool) error {
	if device {
		d := auth.NewDeviceClient(a.Config.BaseURL)
		if err := d.Authorize(ctx); err != nil {
			_ = a.Formatter.Error(err.Error())
			return err
		}
		return nil
	}

	oauthClient := auth.NewOAuthClient(a.Config.BaseURL, a.KeyStore)
	tokens, err := oauthClient.Authenticate(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("authentication failed: %v", err))
		return err
	}

	ttl := time.Duration(tokens.ExpiresIn) * time.Second
	if ttl <= 0 {
		ttl = 3600 * time.Second
	}
	a.TokenCache.SetAccessToken(tokens.AccessToken, ttl)
	a.TokenCache.SetAuthMethod(auth.AuthMethodOAuth)
	a.Client.SetToken(tokens.AccessToken)
	return a.Formatter.Success("Successfully authenticated!")
}

// LoginAPIKey authenticates using an API key.
func (a *App) LoginAPIKey(ctx context.Context, apiKey string) error {
	client := auth.NewAPIKeyClient(a.Config.BaseURL, a.KeyStore)
	profile, err := client.Login(ctx, apiKey)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("API key authentication failed: %v", err))
		return err
	}
	a.Client.SetAPIKey(apiKey)
	a.TokenCache.SetAuthMethod(auth.AuthMethodAPIKey)
	return a.Formatter.Success(fmt.Sprintf("Authenticated as %s (%s) using API key.", profile.Name, profile.Email))
}

// Logout clears stored credentials.
func (a *App) Logout(ctx context.Context) error {
	if err := a.KeyStore.Clear(); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to clear credentials: %v", err))
		return err
	}
	a.Client.SetToken("")
	a.Client.SetAPIKey("")
	return a.Formatter.Success("Logged out successfully.")
}

// Whoami displays the current authenticated user's profile.
func (a *App) Whoami(ctx context.Context) error {
	// Try API key first
	apiKey, _ := a.KeyStore.LoadAPIKey()
	if apiKey != "" {
		client := auth.NewAPIKeyClient(a.Config.BaseURL, a.KeyStore)
		profile, err := client.Validate(ctx, apiKey)
		if err != nil {
			_ = a.Formatter.Error(fmt.Sprintf("failed to get profile: %v", err))
			return err
		}
		if a.Formatter.JSON {
			return a.Formatter.Print(profile)
		}
		_ = a.Formatter.Println(fmt.Sprintf("Name:     %s", profile.Name))
		_ = a.Formatter.Println(fmt.Sprintf("Email:    %s", profile.Email))
		_ = a.Formatter.Println(fmt.Sprintf("Username: %s", profile.Username))
		return a.Formatter.Println("Auth:     API key")
	}

	// OAuth: ensure we have a valid access token
	accessToken, err := a.TokenCache.GetAccessToken(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("not authenticated (run 'mynow login'): %v", err))
		return err
	}
	a.Client.SetToken(accessToken)

	resp, err := a.Client.Get(ctx, "/api/v1/customers", nil)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get profile: %v", err))
		return err
	}
	var profile auth.CustomerProfile
	if err := resp.DecodeJSON(&profile); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to parse profile: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(profile)
	}
	_ = a.Formatter.Println(fmt.Sprintf("Name:     %s", profile.Name))
	_ = a.Formatter.Println(fmt.Sprintf("Email:    %s", profile.Email))
	_ = a.Formatter.Println(fmt.Sprintf("Username: %s", profile.Username))
	return a.Formatter.Println("Auth:     OAuth")
}

// AuthStatus shows the current authentication status.
func (a *App) AuthStatus(ctx context.Context) error {
	method := a.TokenCache.GetAuthMethod()
	if method == "" {
		// Try to detect from stored credentials
		if apiKey, err := a.KeyStore.LoadAPIKey(); err == nil && apiKey != "" {
			method = auth.AuthMethodAPIKey
		} else if _, err := a.KeyStore.LoadRefreshToken(); err == nil {
			method = auth.AuthMethodOAuth
		} else {
			return a.Formatter.Println("Not authenticated. Run 'mynow login' to authenticate.")
		}
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(map[string]interface{}{
			"method":    string(method),
			"expiresAt": a.TokenCache.ExpiresAt().Format(time.RFC3339),
		})
	}
	_ = a.Formatter.Println(fmt.Sprintf("Auth method:   %s", method))
	if method == auth.AuthMethodOAuth && !a.TokenCache.ExpiresAt().IsZero() {
		_ = a.Formatter.Println(fmt.Sprintf("Token expires: %s", a.TokenCache.ExpiresAt().Format("2006-01-02 15:04:05")))
	}
	_ = a.Formatter.Println(fmt.Sprintf("Keyring:       %s", a.Config.Auth.Keyring))
	return nil
}

// AuthRefresh forces a token refresh.
func (a *App) AuthRefresh(ctx context.Context) error {
	if a.TokenCache.GetAuthMethod() == auth.AuthMethodAPIKey {
		return a.Formatter.Info("API key auth does not require token refresh.")
	}
	_, err := a.TokenCache.Refresh(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("token refresh failed: %v", err))
		return err
	}
	return a.Formatter.Success("Token refreshed successfully.")
}

// InboxAdd adds an item to the inbox.
// Inbox items are tasks with null priority — delegate to TaskAdd with no priority.
func (a *App) InboxAdd(ctx context.Context, title string) error {
	return a.TaskAdd(ctx, title, TaskAddOptions{})
}

// InboxList lists inbox items (tasks with null priority).
func (a *App) InboxList(ctx context.Context) error {
	return a.TaskListFull(ctx, TaskListOptions{Priority: "inbox"})
}

// InboxCount prints the count of inbox items.
func (a *App) InboxCount(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	params := api.TaskListParams{Type: "TASK", Priority: ""}
	tasks, err := a.Client.ListTasks(ctx, params)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to count inbox: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(map[string]int{"count": len(tasks)})
	}
	return a.Formatter.Println(fmt.Sprintf("%d", len(tasks)))
}

// InboxProcess interactively walks through inbox items assigning priorities.
func (a *App) InboxProcess(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	params := api.TaskListParams{Priority: ""}
	tasks, err := a.Client.ListTasks(ctx, params)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list inbox: %v", err))
		return err
	}
	if len(tasks) == 0 {
		return a.Formatter.Println("Inbox is empty. 🎉")
	}
	processed := 0
	for _, t := range tasks {
		_ = a.Formatter.Println(fmt.Sprintf("\n> %q", t.Title))
		_ = a.Formatter.Println("  [c]ritical  [o]pportunity  [h]orizon  [p]arking  [s]kip  [d]elete")
		_ = a.Formatter.Printf("  > ")
		var choice string
		if _, err := fmt.Scan(&choice); err != nil {
			break
		}
		switch choice {
		case "c":
			if _, err := a.Client.UpdateTask(ctx, t.ID, api.UpdateTaskRequest{Priority: "CRITICAL"}); err != nil {
				_ = a.Formatter.Error(fmt.Sprintf("failed to update task: %v", err))
			} else {
				_ = a.Formatter.Println(fmt.Sprintf("  ✓ %q → Critical Now", t.Title))
				processed++
			}
		case "o":
			if _, err := a.Client.UpdateTask(ctx, t.ID, api.UpdateTaskRequest{Priority: "OPPORTUNITY_NOW"}); err != nil {
				_ = a.Formatter.Error(fmt.Sprintf("failed to update task: %v", err))
			} else {
				_ = a.Formatter.Println(fmt.Sprintf("  ✓ %q → Opportunity Now", t.Title))
				processed++
			}
		case "h":
			if _, err := a.Client.UpdateTask(ctx, t.ID, api.UpdateTaskRequest{Priority: "OVER_THE_HORIZON"}); err != nil {
				_ = a.Formatter.Error(fmt.Sprintf("failed to update task: %v", err))
			} else {
				_ = a.Formatter.Println(fmt.Sprintf("  ✓ %q → Over The Horizon", t.Title))
				processed++
			}
		case "p":
			if _, err := a.Client.UpdateTask(ctx, t.ID, api.UpdateTaskRequest{Priority: "PARKING_LOT"}); err != nil {
				_ = a.Formatter.Error(fmt.Sprintf("failed to update task: %v", err))
			} else {
				_ = a.Formatter.Println(fmt.Sprintf("  ✓ %q → Parking Lot", t.Title))
				processed++
			}
		case "d":
			if err := a.Client.DeleteTask(ctx, t.ID, false); err != nil {
				_ = a.Formatter.Error(fmt.Sprintf("failed to delete task: %v", err))
			} else {
				_ = a.Formatter.Println(fmt.Sprintf("  ✗ %q deleted", t.Title))
				processed++
			}
		default:
			_ = a.Formatter.Println("  → skipped")
		}
	}
	return a.Formatter.Println(fmt.Sprintf("\nProcessed %d of %d items. %d remaining.", processed, len(tasks), len(tasks)-processed))
}

// InboxClear archives all inbox items.
func (a *App) InboxClear(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	params := api.TaskListParams{Priority: ""}
	tasks, err := a.Client.ListTasks(ctx, params)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list inbox: %v", err))
		return err
	}
	ids := make([]string, 0, len(tasks))
	for _, t := range tasks {
		ids = append(ids, t.ID)
	}
	if len(ids) == 0 {
		return a.Formatter.Println("Inbox is already empty.")
	}
	// Archive by deleting each individually (BatchUpdate doesn't support archived field directly)
	var failed int
	for _, id := range ids {
		if err := a.Client.DeleteTask(ctx, id, false); err != nil {
			failed++
			_ = a.Formatter.Error(fmt.Sprintf("failed to clear task %s: %v", id, err))
		}
	}
	cleared := len(ids) - failed
	return a.Formatter.Success(fmt.Sprintf("Cleared %d of %d inbox items.", cleared, len(ids)))
}

// NowList lists current focus items (CRITICAL + OPPORTUNITY_NOW tasks for today).
func (a *App) NowList(ctx context.Context) error {
	return a.TaskListFull(ctx, TaskListOptions{Today: true})
}

// NowFocus shows current CRITICAL tasks (show mode).
func (a *App) NowFocus(ctx context.Context) error {
	return a.TaskListFull(ctx, TaskListOptions{Priority: "critical"})
}

// NowFocusSet sets a task to CRITICAL priority (making it a focus task).
func (a *App) NowFocusSet(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	task, err := a.Client.UpdateTask(ctx, id, api.UpdateTaskRequest{Priority: "CRITICAL"})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to set focus: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	return a.Formatter.Success(fmt.Sprintf("● Focus set: %s", task.Title))
}

// NowFocusClear moves CRITICAL tasks back to OPPORTUNITY_NOW (clearing focus).
func (a *App) NowFocusClear(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	params := api.TaskListParams{Priority: "CRITICAL"}
	tasks, err := a.Client.ListTasks(ctx, params)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get focus tasks: %v", err))
		return err
	}
	ids := make([]string, 0, len(tasks))
	for _, t := range tasks {
		ids = append(ids, t.ID)
	}
	if len(ids) == 0 {
		return a.Formatter.Println("No focus tasks to clear.")
	}
	batchReq := api.BatchUpdateRequest{IDs: ids, Updates: api.UpdateTaskRequest{Priority: "OPPORTUNITY_NOW"}}
	if _, err := a.Client.BatchUpdateTasks(ctx, batchReq); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to clear focus: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Cleared focus for %d task(s).", len(ids)))
}

// NowComplete completes a task from the now view.
func (a *App) NowComplete(ctx context.Context, id string) error {
	return a.TaskComplete(ctx, id)
}

// NowSnooze snoozes a task from the now view.
func (a *App) NowSnooze(ctx context.Context, id, date string, days int) error {
	return a.TaskSnoozeTask(ctx, id, TaskSnoozeOpt{Date: date, Days: days})
}

// ReviewDaily runs the daily review — walks through overdue/uncompleted tasks interactively.
func (a *App) ReviewDaily(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}

	now := time.Now()
	dayName := now.Format("Monday, January 2")
	_ = a.Formatter.Println(fmt.Sprintf("DAILY REVIEW — %s\n", dayName))

	// Fetch overdue and today's incomplete tasks
	params := api.TaskListParams{}
	tasks, err := a.Client.ListTasks(ctx, params)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list tasks: %v", err))
		return err
	}

	today := now.Format("2006-01-02")
	var reviewTasks []api.UnifiedTask
	for _, t := range tasks {
		if t.IsCompleted || t.IsArchived {
			continue
		}
		if t.StartDate != "" && t.StartDate <= today {
			reviewTasks = append(reviewTasks, t)
		}
	}

	if len(reviewTasks) == 0 {
		return a.Formatter.Println("No tasks to review. Great work!")
	}

	_ = a.Formatter.Println(fmt.Sprintf("Checking in on %d tasks...\n", len(reviewTasks)))

	tomorrow := now.AddDate(0, 0, 1).Format("2006-01-02")
	processed := 0
	for _, t := range reviewTasks {
		overdue := t.StartDate < today
		dateLabel := ""
		if overdue {
			dateLabel = fmt.Sprintf("  Due: %s (overdue)", t.StartDate)
		} else {
			dateLabel = fmt.Sprintf("  Due: %s", t.StartDate)
		}
		symbol := output.PriorityColored(t.PriorityString(), a.Formatter.NoColor)
		_ = a.Formatter.Println(fmt.Sprintf("> %s  [%s]%s", t.Title, symbol, dateLabel))
		_ = a.Formatter.Println("  [c]omplete  [s]nooze  [r]eschedule  [k]eep  [d]elete")
		_ = a.Formatter.Printf("  > ")

		var choice string
		if _, err := fmt.Scan(&choice); err != nil {
			break
		}
		switch choice {
		case "c":
			if _, err := a.Client.CompleteTask(ctx, t.ID); err != nil {
				_ = a.Formatter.Error(fmt.Sprintf("failed to complete task: %v", err))
			} else {
				_ = a.Formatter.Println(fmt.Sprintf("  Completed: %s", t.Title))
				processed++
			}
		case "s":
			if _, err := a.Client.UpdateTask(ctx, t.ID, api.UpdateTaskRequest{StartDate: tomorrow}); err != nil {
				_ = a.Formatter.Error(fmt.Sprintf("failed to snooze task: %v", err))
			} else {
				_ = a.Formatter.Println("  Snoozed to tomorrow.")
				processed++
			}
		case "r":
			_ = a.Formatter.Printf("  New date (YYYY-MM-DD): ")
			var newDate string
			if _, err := fmt.Scan(&newDate); err == nil && newDate != "" {
				if _, err := a.Client.UpdateTask(ctx, t.ID, api.UpdateTaskRequest{StartDate: newDate}); err != nil {
					_ = a.Formatter.Error(fmt.Sprintf("failed to reschedule task: %v", err))
				} else {
					_ = a.Formatter.Println(fmt.Sprintf("  Rescheduled to %s.", newDate))
					processed++
				}
			}
		case "d":
			if err := a.Client.DeleteTask(ctx, t.ID, false); err != nil {
				_ = a.Formatter.Error(fmt.Sprintf("failed to delete task: %v", err))
			} else {
				_ = a.Formatter.Println(fmt.Sprintf("  Deleted: %s", t.Title))
				processed++
			}
		default:
			_ = a.Formatter.Println("  Kept.")
		}
	}

	return a.Formatter.Println(fmt.Sprintf("\nReview complete. %d tasks processed.", processed))
}

// ReviewWeekly runs the weekly review — shows the week's completions and plans next week.
func (a *App) ReviewWeekly(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}

	now := time.Now()
	// Week boundaries: Monday–Sunday of the current week
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7 in ISO
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
	weekEnd := weekStart.AddDate(0, 0, 6)
	weekStartStr := weekStart.Format("Jan 2")
	weekEndStr := weekEnd.Format("Jan 2, 2006")
	_ = a.Formatter.Println(fmt.Sprintf("WEEKLY REVIEW — Week of %s-%s\n", weekStartStr, weekEndStr))

	// Fetch completed tasks
	completedParams := api.TaskListParams{IsCompleted: true}
	completedTasks, err := a.Client.ListTasks(ctx, completedParams)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list completed tasks: %v", err))
		return err
	}

	// Fetch all active tasks
	activeTasks, err := a.Client.ListTasks(ctx, api.TaskListParams{})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list tasks: %v", err))
		return err
	}

	_ = a.Formatter.Println(fmt.Sprintf("COMPLETED THIS WEEK: %d tasks\n", len(completedTasks)))

	// Separate overdue vs upcoming incomplete tasks
	today := now.Format("2006-01-02")
	nextWeekStart := now.AddDate(0, 0, 7-weekday+1).Format("2006-01-02")
	nextWeekEnd := now.AddDate(0, 0, 7-weekday+7).Format("2006-01-02")

	var overdueTasks, nextWeekTasks []api.UnifiedTask
	for _, t := range activeTasks {
		if t.IsCompleted || t.IsArchived {
			continue
		}
		if t.StartDate != "" && t.StartDate < today {
			overdueTasks = append(overdueTasks, t)
		} else if t.StartDate >= nextWeekStart && t.StartDate <= nextWeekEnd {
			nextWeekTasks = append(nextWeekTasks, t)
		}
	}

	_ = a.Formatter.Println("INCOMPLETE:")
	if len(overdueTasks) == 0 {
		_ = a.Formatter.Println("  (none)")
	} else {
		for _, t := range overdueTasks {
			symbol := output.PriorityColored(t.PriorityString(), a.Formatter.NoColor)
			_ = a.Formatter.Println(fmt.Sprintf("  %s %s (overdue)", symbol, t.Title))
		}
	}

	_ = a.Formatter.Println("\nNEXT WEEK:")
	if len(nextWeekTasks) == 0 {
		_ = a.Formatter.Println("  (no tasks scheduled)")
	} else {
		for _, t := range nextWeekTasks {
			symbol := output.PriorityColored(t.PriorityString(), a.Formatter.NoColor)
			_ = a.Formatter.Println(fmt.Sprintf("  %s %s  [%s]", symbol, t.Title, t.StartDate))
		}
	}

	return a.Formatter.Println("\nReview complete.")
}

// RunTUI is deprecated — TUI is launched directly via tui.Run() from the cmd layer.
// Kept only for interface compatibility if any external code references it.
func (a *App) RunTUI(ctx context.Context) error {
	return fmt.Errorf("use tui.Run(app) directly; this method is deprecated")
}

// derefBool safely dereferences a *bool pointer, returning false if nil.
func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// ConfigShow prints the resolved configuration (secrets redacted).
func (a *App) ConfigShow(ctx context.Context) error {
	cfg := a.Config
	if a.Formatter.JSON {
		type jsonCfg struct {
			ConfigFile string               `json:"config_file"`
			API        config.APIConfig     `json:"api"`
			Auth       config.AuthConfig    `json:"auth"`
			Display    config.DisplayConfig `json:"display"`
			TUI        config.TUIConfig     `json:"tui"`
			Defaults   config.DefaultsConfig `json:"defaults"`
			APIKey     string               `json:"api_key,omitempty"`
		}
		out := jsonCfg{
			ConfigFile: cfg.ConfigFile,
			API:        cfg.API,
			Auth:       cfg.Auth,
			Display:    cfg.Display,
			TUI:        cfg.TUI,
			Defaults:   cfg.Defaults,
		}
		if cfg.APIKey != "" {
			out.APIKey = "***redacted***"
		}
		return a.Formatter.Print(out)
	}
	lines := []string{
		fmt.Sprintf("config file:                   %s", cfg.ConfigFile),
		fmt.Sprintf("api.url:                       %s", cfg.API.URL),
		fmt.Sprintf("api.timeout:                   %s", cfg.API.Timeout),
		fmt.Sprintf("api.retries:                   %d", cfg.API.Retries),
		fmt.Sprintf("auth.method:                   %s", cfg.Auth.Method),
		fmt.Sprintf("auth.keyring:                  %s", cfg.Auth.Keyring),
		fmt.Sprintf("display.color:                 %s", cfg.Display.Color),
		fmt.Sprintf("display.date_format:           %s", cfg.Display.DateFormat),
		fmt.Sprintf("display.time_format:           %s", cfg.Display.TimeFormat),
		fmt.Sprintf("display.default_output:        %s", cfg.Display.DefaultOutput),
		fmt.Sprintf("tui.theme:                     %s", cfg.TUI.Theme),
		fmt.Sprintf("tui.refresh_interval:          %s", cfg.TUI.RefreshInterval),
		fmt.Sprintf("tui.vim_keys:                  %v", derefBool(cfg.TUI.VimKeys)),
		fmt.Sprintf("tui.mouse:                     %v", derefBool(cfg.TUI.Mouse)),
		fmt.Sprintf("tui.animations:                %v", derefBool(cfg.TUI.Animations)),
		fmt.Sprintf("defaults.priority:             %s", cfg.Defaults.Priority),
		fmt.Sprintf("defaults.task_type:            %s", cfg.Defaults.TaskType),
		fmt.Sprintf("defaults.calendar_days:        %d", cfg.Defaults.CalendarDays),
		fmt.Sprintf("defaults.habit_schedule_days:  %d", cfg.Defaults.HabitScheduleDays),
	}
	if cfg.APIKey != "" {
		lines = append(lines, "api_key:                       ***redacted***")
	}
	for _, line := range lines {
		if err := a.Formatter.Println(line); err != nil {
			return err
		}
	}
	return nil
}

// ConfigGet prints the value of a single config key.
func (a *App) ConfigGet(ctx context.Context, key string) error {
	val, err := config.GetValue(a.Config, key)
	if err != nil {
		_ = a.Formatter.Error(err.Error())
		return err
	}
	return a.Formatter.Println(val)
}

// ConfigSet sets a config key and persists it to the config file.
func (a *App) ConfigSet(ctx context.Context, key, value string) error {
	if err := config.SetValue(a.Config, key, value); err != nil {
		_ = a.Formatter.Error(err.Error())
		return err
	}
	if err := config.Save(a.Config); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to save config: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Set %s = %s", key, value))
}

// ConfigReset removes the config file, reverting to defaults.
func (a *App) ConfigReset(ctx context.Context) error {
	if err := config.Reset(a.Config); err != nil {
		_ = a.Formatter.Error(err.Error())
		return err
	}
	return a.Formatter.Success("Configuration reset to defaults.")
}

// ConfigPath prints the path to the config file.
func (a *App) ConfigPath(ctx context.Context) error {
	return a.Formatter.Println(a.Config.ConfigFile)
}
