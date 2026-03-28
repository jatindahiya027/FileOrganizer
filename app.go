package main

import (
	"context"
	"file-organizer/organizer"
	goruntime "runtime"
	"strings"
	"sync"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the Wails application struct
type App struct {
	ctx       context.Context
	mu        sync.Mutex
	org       *organizer.Organizer
	isRunning bool
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called at application shutdown
func (a *App) shutdown(ctx context.Context) {
	a.Cancel()
}

// SelectDirectory opens a native folder picker and returns selected path
func (a *App) SelectDirectory() string {
	path, err := wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "Select folder to organize",
	})
	if err != nil {
		return ""
	}
	return path
}

// StartOrganize begins organizing files in the given path
func (a *App) StartOrganize(path string, dryRun bool, threads int, excludeExts string) organizer.Summary {
	a.mu.Lock()
	if a.isRunning {
		a.mu.Unlock()
		return organizer.Summary{}
	}
	a.isRunning = true
	a.mu.Unlock()

	defer func() {
		a.mu.Lock()
		a.isRunning = false
		a.mu.Unlock()
		// organizer already emits final progress; just mark not running
	}()

	// Parse excluded extensions
	var excludeList []string
	if excludeExts != "" {
		for _, e := range strings.Split(excludeExts, ",") {
			e = strings.TrimSpace(e)
			if e != "" {
				if !strings.HasPrefix(e, ".") {
					e = "." + e
				}
				excludeList = append(excludeList, strings.ToLower(e))
			}
		}
	}

	cfg := organizer.OrganizerConfig{
		Path:        path,
		DryRun:      dryRun,
		Threads:     threads,
		ExcludeExts: excludeList,
	}

	org := organizer.New(cfg,
		func(entry organizer.LogEntry) {
			wailsruntime.EventsEmit(a.ctx, "log", entry)
		},
		func(prog organizer.ProgressEvent) {
			wailsruntime.EventsEmit(a.ctx, "progress", prog)
		},
	)

	a.mu.Lock()
	a.org = org
	a.mu.Unlock()

	return org.Run(context.Background())
}

// Cancel stops the current organization run
func (a *App) Cancel() {
	a.mu.Lock()
	org := a.org
	a.mu.Unlock()
	if org != nil {
		org.Cancel()
	}
}

// GetProgress returns the current progress snapshot
func (a *App) GetProgress() organizer.ProgressEvent {
	a.mu.Lock()
	org := a.org
	a.mu.Unlock()
	if org != nil {
		return org.GetProgress()
	}
	return organizer.ProgressEvent{}
}

// IsRunning returns whether an organization is currently in progress
func (a *App) IsRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.isRunning
}

// GetCPUCount returns the number of logical CPUs for thread slider default
func (a *App) GetCPUCount() int {
	return goruntime.NumCPU()
}

// GetCategories returns the full extension→category map
func (a *App) GetCategories() map[string][]string {
	reverse := make(map[string][]string)
	for ext, cat := range organizer.ExtensionMap {
		reverse[cat] = append(reverse[cat], ext)
	}
	return reverse
}

// GetCategoryColors returns category→color map for the frontend
func (a *App) GetCategoryColors() map[string]string {
	return organizer.CategoryColor
}
