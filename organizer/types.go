package organizer

import "time"

// FileTask represents a single file to be moved
type FileTask struct {
	SrcPath  string
	DstDir   string
	Filename string
	Category string
	DryRun   bool
}

// MoveResult is the result of a single file move
type MoveResult struct {
	Filename string    `json:"filename"`
	SrcPath  string    `json:"srcPath"`
	DstPath  string    `json:"dstPath"`
	Category string    `json:"category"`
	Success  bool      `json:"success"`
	Skipped  bool      `json:"skipped"`
	DryRun   bool      `json:"dryRun"`
	Error    string    `json:"error,omitempty"`
	Time     time.Time `json:"time"`
}

// LogEntry is a structured log line emitted in real time
type LogEntry struct {
	Level    string    `json:"level"` // "info" | "success" | "warn" | "error"
	Message  string    `json:"message"`
	Filename string    `json:"filename,omitempty"`
	Category string    `json:"category,omitempty"`
	Time     time.Time `json:"time"`
}

// ProgressEvent is emitted frequently during organization
type ProgressEvent struct {
	Total          int            `json:"total"`
	Processed      int            `json:"processed"`
	Moved          int            `json:"moved"`
	Skipped        int            `json:"skipped"`
	Errors         int            `json:"errors"`
	CurrentFile    string         `json:"currentFile"`
	PercentDone    float64        `json:"percentDone"`
	CategoryCounts map[string]int `json:"categoryCounts"`
	Running        bool           `json:"running"`
}

// Summary returned after organization completes
type Summary struct {
	Total          int            `json:"total"`
	Moved          int            `json:"moved"`
	Skipped        int            `json:"skipped"`
	Errors         int            `json:"errors"`
	ElapsedSeconds float64        `json:"elapsedSeconds"`
	CategoryCounts map[string]int `json:"categoryCounts"`
	LogPath        string         `json:"logPath"`
	DryRun         bool           `json:"dryRun"`
}

// OrganizerConfig holds runtime config
type OrganizerConfig struct {
	Path        string
	DryRun      bool
	Threads     int
	ExcludeExts []string
}
