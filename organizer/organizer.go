package organizer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// SystemSkipFiles are files that should never be moved
var SystemSkipFiles = map[string]bool{
	"desktop.ini":          true,
	"thumbs.db":            true,
	".ds_store":            true,
	".spotlight-v100":      true,
	"organizer.log":        true,
	"file-organizer":       true,
	"file-organizer.exe":   true,
	"file organizer":       true,
	"file organizer.exe":   true,
}

// Organizer runs the file organization with a worker pool
type Organizer struct {
	config      OrganizerConfig
	cancelFunc  context.CancelFunc
	mu          sync.Mutex
	logEntries  []LogEntry
	progress    ProgressEvent
	onLog       func(LogEntry)
	onProgress  func(ProgressEvent)
}

// New creates an Organizer
func New(cfg OrganizerConfig, onLog func(LogEntry), onProgress func(ProgressEvent)) *Organizer {
	return &Organizer{
		config:     cfg,
		onLog:      onLog,
		onProgress: onProgress,
	}
}

// Run starts the organization process and returns a Summary
func (o *Organizer) Run(ctx context.Context) Summary {
	start := time.Now()
	ctx, cancel := context.WithCancel(ctx)
	o.cancelFunc = cancel
	defer cancel()

	o.emit(LogEntry{Level: "info", Message: fmt.Sprintf("Scanning directory: %s", o.config.Path)})

	// Scan only the top-level directory (no subdirectories)
	files, err := o.scanDirectory()
	if err != nil {
		o.emit(LogEntry{Level: "error", Message: fmt.Sprintf("Failed to scan directory: %s", err)})
		return Summary{}
	}

	total := len(files)
	o.emit(LogEntry{Level: "info", Message: fmt.Sprintf("Found %d files to organize", total)})

	if total == 0 {
		o.emit(LogEntry{Level: "warn", Message: "No files found to organize"})
		return Summary{Total: 0, ElapsedSeconds: time.Since(start).Seconds()}
	}

	// Initialize progress
	o.mu.Lock()
	o.progress = ProgressEvent{
		Total:          total,
		CategoryCounts: make(map[string]int),
		Running:        true,
	}
	o.mu.Unlock()

	// Build tasks
	tasks := make(chan FileTask, total)
	results := make(chan MoveResult, total)

	// Enqueue tasks
	for _, f := range files {
		tasks <- f
	}
	close(tasks)

	// Determine thread count
	threads := o.config.Threads
	if threads <= 0 {
		threads = runtime.NumCPU() * 2
	}
	if threads > 32 {
		threads = 32
	}
	if threads > total {
		threads = total
	}

	o.emit(LogEntry{Level: "info", Message: fmt.Sprintf("Starting %d worker threads", threads)})

	// Create destination folders upfront
	if !o.config.DryRun {
		o.createCategoryFolders(files)
	}

	// Launch workers
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for task := range tasks {
				select {
				case <-ctx.Done():
					results <- MoveResult{
						Filename: task.Filename,
						SrcPath:  task.SrcPath,
						Category: task.Category,
						Skipped:  true,
						Error:    "cancelled",
					}
					return
				default:
					result := o.processTask(task)
					results <- result
				}
			}
		}(i)
	}

	// Close results when all workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results and update progress
	moved, skipped, errors := 0, 0, 0
	catCounts := make(map[string]int)
	var logLines []string

	processed := 0
	for r := range results {
		processed++
		logLine, _ := json.Marshal(map[string]interface{}{
			"time":     time.Now().Format(time.RFC3339),
			"action":   resultAction(r),
			"filename": r.Filename,
			"src":      r.SrcPath,
			"dst":      r.DstPath,
			"category": r.Category,
			"error":    r.Error,
		})
		logLines = append(logLines, string(logLine))

		if r.Success {
			moved++
			catCounts[r.Category]++
			o.emit(LogEntry{
				Level:    "success",
				Message:  fmt.Sprintf("→ %s  [%s]", r.Filename, r.Category),
				Filename: r.Filename,
				Category: r.Category,
			})
		} else if r.Skipped {
			skipped++
			if r.Error != "cancelled" {
				o.emit(LogEntry{
					Level:    "warn",
					Message:  fmt.Sprintf("⊘ Skipped: %s", r.Filename),
					Filename: r.Filename,
				})
			}
		} else {
			errors++
			o.emit(LogEntry{
				Level:    "error",
				Message:  fmt.Sprintf("✕ Error moving %s: %s", r.Filename, r.Error),
				Filename: r.Filename,
			})
		}

		pct := float64(processed) / float64(total) * 100
		o.mu.Lock()
		o.progress = ProgressEvent{
			Total:          total,
			Processed:      processed,
			Moved:          moved,
			Skipped:        skipped,
			Errors:         errors,
			CurrentFile:    r.Filename,
			PercentDone:    pct,
			CategoryCounts: catCounts,
			Running:        processed < total,
		}
		prog := o.progress
		o.mu.Unlock()

		if o.onProgress != nil {
			o.onProgress(prog)
		}
	}

	elapsed := time.Since(start).Seconds()

	// Write log file
	logPath := filepath.Join(o.config.Path, "organizer.log")
	if !o.config.DryRun {
		logContent := strings.Join(logLines, "\n") + "\n"
		os.WriteFile(logPath, []byte(logContent), 0644)
	}

	mode := ""
	if o.config.DryRun {
		mode = " (DRY RUN)"
	}
	o.emit(LogEntry{
		Level:   "info",
		Message: fmt.Sprintf("✓ Done%s — Moved: %d  Skipped: %d  Errors: %d  Time: %.2fs", mode, moved, skipped, errors, elapsed),
	})

	return Summary{
		Total:          total,
		Moved:          moved,
		Skipped:        skipped,
		Errors:         errors,
		ElapsedSeconds: elapsed,
		CategoryCounts: catCounts,
		LogPath:        logPath,
		DryRun:         o.config.DryRun,
	}
}

// Cancel stops the worker pool
func (o *Organizer) Cancel() {
	if o.cancelFunc != nil {
		o.cancelFunc()
		o.emit(LogEntry{Level: "warn", Message: "⚡ Organization cancelled by user"})
	}
}

// GetProgress returns current progress snapshot
func (o *Organizer) GetProgress() ProgressEvent {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.progress
}

// scanDirectory lists only top-level files (no subdirs)
func (o *Organizer) scanDirectory() ([]FileTask, error) {
	entries, err := os.ReadDir(o.config.Path)
	if err != nil {
		return nil, err
	}

	excludeSet := make(map[string]bool)
	for _, e := range o.config.ExcludeExts {
		excludeSet[strings.ToLower(e)] = true
	}

	var tasks []FileTask
	for _, entry := range entries {
		if entry.IsDir() {
			continue // never recurse — only flat scan
		}

		name := entry.Name()
		nameLower := strings.ToLower(name)

		// Skip system files
		if SystemSkipFiles[nameLower] {
			continue
		}

		// Skip hidden files
		if strings.HasPrefix(name, ".") {
			continue
		}

		ext := strings.ToLower(filepath.Ext(name))

		// Skip excluded extensions
		if excludeSet[ext] {
			continue
		}

		category := GetCategory(ext)
		dstDir := filepath.Join(o.config.Path, category)

		tasks = append(tasks, FileTask{
			SrcPath:  filepath.Join(o.config.Path, name),
			DstDir:   dstDir,
			Filename: name,
			Category: category,
			DryRun:   o.config.DryRun,
		})
	}
	return tasks, nil
}

// createCategoryFolders pre-creates all needed destination folders
func (o *Organizer) createCategoryFolders(tasks []FileTask) {
	seen := make(map[string]bool)
	for _, t := range tasks {
		if !seen[t.DstDir] {
			seen[t.DstDir] = true
			if err := os.MkdirAll(t.DstDir, 0755); err == nil {
				o.emit(LogEntry{Level: "info", Message: fmt.Sprintf("📁 Created folder: %s", filepath.Base(t.DstDir))})
			}
		}
	}
}

// processTask performs the actual move for one file
func (o *Organizer) processTask(task FileTask) MoveResult {
	result := MoveResult{
		Filename: task.Filename,
		SrcPath:  task.SrcPath,
		Category: task.Category,
		DryRun:   task.DryRun,
		Time:     time.Now(),
	}

	if task.DryRun {
		result.DstPath = filepath.Join(task.DstDir, task.Filename)
		result.Success = true
		return result
	}

	dstPath, err := MoveFile(task.SrcPath, task.DstDir, task.Filename)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.DstPath = dstPath
	result.Success = true
	return result
}

// emit sends a log entry to the callback
func (o *Organizer) emit(entry LogEntry) {
	entry.Time = time.Now()
	o.mu.Lock()
	o.logEntries = append(o.logEntries, entry)
	o.mu.Unlock()
	if o.onLog != nil {
		o.onLog(entry)
	}
}

func resultAction(r MoveResult) string {
	if r.Success {
		if r.DryRun {
			return "dry-run"
		}
		return "moved"
	}
	if r.Skipped {
		return "skipped"
	}
	return "error"
}
