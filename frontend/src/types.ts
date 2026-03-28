// Type declarations for Wails-generated bindings
// Wails injects these globals automatically — no npm install needed

export interface LogEntry {
  level: 'info' | 'success' | 'warn' | 'error'
  message: string
  filename?: string
  category?: string
  time: string
}

export interface ProgressEvent {
  total: number
  processed: number
  moved: number
  skipped: number
  errors: number
  currentFile: string
  percentDone: number
  categoryCounts: Record<string, number>
  running: boolean
}

export interface Summary {
  total: number
  moved: number
  skipped: number
  errors: number
  elapsedSeconds: number
  categoryCounts: Record<string, number>
  logPath: string
  dryRun: boolean
}

// Wails injects window.go and window.runtime at runtime
declare global {
  interface Window {
    go: {
      main: {
        App: {
          SelectDirectory(): Promise<string>
          StartOrganize(path: string, dryRun: boolean, threads: number, excludeExts: string): Promise<Summary>
          Cancel(): Promise<void>
          GetProgress(): Promise<ProgressEvent>
          IsRunning(): Promise<boolean>
          GetCPUCount(): Promise<number>
          GetCategories(): Promise<Record<string, string[]>>
          GetCategoryColors(): Promise<Record<string, string>>
        }
      }
    }
    runtime: {
      EventsOn(event: string, callback: (...args: any[]) => void): void
      EventsOff(event: string): void
    }
  }
}
