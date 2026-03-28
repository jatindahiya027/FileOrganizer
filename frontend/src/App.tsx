import { useState, useEffect, useRef, useCallback } from 'react'
import type { LogEntry, ProgressEvent, Summary } from './types'

// ── helpers ────────────────────────────────────────────────────────────────
const wails = () => window.go?.main?.App
const rt = () => window.runtime

const CATEGORY_COLORS: Record<string, string> = {
  Images: '#60a5fa', Videos: '#f472b6', Audio: '#a78bfa',
  Documents: '#34d399', Spreadsheets: '#6ee7b7', Presentations: '#fbbf24',
  Ebooks: '#fb923c', '3D_Files': '#f87171', Executables: '#e879f9',
  Archives: '#94a3b8', Disk_Images: '#475569', Code: '#38bdf8',
  Fonts: '#c084fc', Database: '#4ade80', Others: '#64748b',
}

const LOG_COLORS: Record<string, string> = {
  success: '#00e676', warn: '#ffb300', error: '#ff5252', info: '#00d4ff',
}
const LOG_ICONS: Record<string, string> = {
  success: '✓', warn: '⊘', error: '✕', info: '·',
}

function formatTime(d: Date) {
  return d.toTimeString().slice(0, 8)
}

// ── sub-components ─────────────────────────────────────────────────────────

function TitleBar() {
  return (
    <div className="drag-region h-10 flex items-center px-5 border-b border-white/5 shrink-0">
      <div className="flex items-center gap-2 no-drag">
        <div className="w-2 h-2 rounded-full bg-cyan-400 animate-pulse-dot" />
        <span className="font-display text-lg tracking-widest text-cyan-400 text-glow-cyan">
          FILE ORGANIZER
        </span>
      </div>
      <div className="ml-auto text-xs font-mono text-slate-600 no-drag">v1.0.0</div>
    </div>
  )
}

function DirectoryPicker({ path, onSelect }: { path: string; onSelect: (p: string) => void }) {
  const handleClick = async () => {
    const p = await wails()?.SelectDirectory()
    if (p) onSelect(p)
  }
  return (
    <div className="space-y-1.5">
      <label className="text-xs font-mono text-slate-500 uppercase tracking-widest">Target Directory</label>
      <div
        onClick={handleClick}
        className="no-drag flex items-center gap-3 px-3 py-2.5 rounded-lg border cursor-pointer transition-all
          border-slate-700 hover:border-cyan-500/50 hover:bg-cyan-500/5 group"
      >
        <svg className="w-4 h-4 text-cyan-500 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
            d="M3 7a2 2 0 012-2h4l2 2h8a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V7z" />
        </svg>
        {path
          ? <span className="text-xs font-mono text-cyan-300 truncate flex-1">{path}</span>
          : <span className="text-xs font-mono text-slate-600 flex-1">Click to select folder…</span>
        }
        <svg className="w-3.5 h-3.5 text-slate-600 group-hover:text-cyan-500 transition-colors shrink-0"
          fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
        </svg>
      </div>
    </div>
  )
}

function ThreadSlider({ value, onChange, max }: { value: number; onChange: (v: number) => void; max: number }) {
  const sliderMax = Math.max(max, 16)
  return (
    <div className="space-y-1.5">
      <div className="flex justify-between items-center">
        <label className="text-xs font-mono text-slate-500 uppercase tracking-widest">Worker Threads</label>
        <span className="text-xs font-mono text-cyan-400 font-semibold">{value}</span>
      </div>
      <input
        type="range" min={1} max={sliderMax} value={value}
        onChange={e => onChange(Number(e.target.value))}
        className="no-drag w-full h-1 rounded-full appearance-none cursor-pointer"
        style={{
          background: `linear-gradient(to right, #00d4ff ${(value / sliderMax) * 100}%, #1a2333 0%)`
        }}
      />
      <div className="flex justify-between text-xs font-mono text-slate-700">
        <span>1</span><span>Auto</span><span>{sliderMax}</span>
      </div>
    </div>
  )
}

function ToggleSwitch({ checked, onChange, label }: { checked: boolean; onChange: (v: boolean) => void; label: string }) {
  return (
    <div className="no-drag flex items-center justify-between">
      <span className="text-xs font-mono text-slate-400">{label}</span>
      <button
        onClick={() => onChange(!checked)}
        className={`relative w-9 h-5 rounded-full transition-colors duration-200 ${checked ? 'bg-amber-500' : 'bg-slate-700'}`}
      >
        <span className={`absolute top-0.5 left-0.5 w-4 h-4 rounded-full bg-white transition-transform duration-200 ${checked ? 'translate-x-4' : ''}`} />
      </button>
    </div>
  )
}

function ProgressBar({ percent, running, done }: { percent: number; running: boolean; done: boolean }) {
  const color = done ? '#00e676' : '#00d4ff'
  return (
    <div className="space-y-1">
      <div className="flex justify-between text-xs font-mono text-slate-500">
        <span>{running ? 'Processing…' : done ? 'Complete' : 'Ready'}</span>
        <span style={{ color }}>{percent.toFixed(1)}%</span>
      </div>
      <div className="h-1.5 bg-surface-3 rounded-full overflow-hidden">
        <div
          className="h-full rounded-full transition-all duration-300"
          style={{
            width: `${percent}%`,
            background: running
              ? `linear-gradient(90deg, ${color}, ${color}88, ${color})`
              : color,
            backgroundSize: running ? '200% 100%' : undefined,
            animation: running ? 'shimmer 1.5s linear infinite' : undefined,
            boxShadow: `0 0 8px ${color}66`,
          }}
        />
      </div>
    </div>
  )
}

function StatPill({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <div className="flex flex-col items-center gap-0.5 px-3 py-2 rounded-lg bg-surface-2 border border-white/5">
      <span className="text-base font-mono font-bold" style={{ color }}>{value}</span>
      <span className="text-xs font-mono text-slate-600">{label}</span>
    </div>
  )
}

function LogLine({ entry }: { entry: LogEntry & { _ts: number } }) {
  const color = LOG_COLORS[entry.level] || '#94a3b8'
  const icon = LOG_ICONS[entry.level] || '·'
  return (
    <div className="flex items-start gap-2 text-xs font-mono py-0.5 animate-slide-in hover:bg-white/2 px-1 rounded">
      <span className="shrink-0 text-slate-700">{formatTime(new Date(entry._ts))}</span>
      <span className="shrink-0 w-3" style={{ color }}>{icon}</span>
      {entry.category && (
        <span className="shrink-0 px-1 rounded text-xs leading-tight"
          style={{ color, background: color + '1a' }}>
          {entry.category}
        </span>
      )}
      <span style={{ color: entry.level === 'info' ? '#64748b' : '#c9d4e8' }}>
        {entry.message}
      </span>
    </div>
  )
}

function CategoryBar({ name, count, total, color }: { name: string; count: number; total: number; color: string }) {
  const pct = total > 0 ? (count / total) * 100 : 0
  return (
    <div className="space-y-1">
      <div className="flex justify-between text-xs font-mono">
        <span className="text-slate-400">{name.replace('_', ' ')}</span>
        <span style={{ color }}>{count}</span>
      </div>
      <div className="h-1 bg-surface-3 rounded-full overflow-hidden">
        <div className="h-full rounded-full transition-all duration-500"
          style={{ width: `${pct}%`, background: color, boxShadow: `0 0 4px ${color}66` }} />
      </div>
    </div>
  )
}

// ── Main App ───────────────────────────────────────────────────────────────

const EMPTY_PROGRESS: ProgressEvent = {
  total: 0, processed: 0, moved: 0, skipped: 0, errors: 0,
  currentFile: '', percentDone: 0, categoryCounts: {}, running: false,
}

export default function App() {
  const [path, setPath] = useState('')
  const [dryRun, setDryRun] = useState(false)
  const [threads, setThreads] = useState(4)
  const [cpuCount, setCpuCount] = useState(4)
  const [excludeExts, setExcludeExts] = useState('')
  const [running, setRunning] = useState(false)
  const [done, setDone] = useState(false)
  // Keep a stable "last good" progress so completion doesn't blank the UI
  const [progress, setProgress] = useState<ProgressEvent>(EMPTY_PROGRESS)
  const [summary, setSummary] = useState<Summary | null>(null)
  const [logs, setLogs] = useState<Array<LogEntry & { _ts: number }>>([])
  const logRef = useRef<HTMLDivElement>(null)
  const [autoScroll, setAutoScroll] = useState(true)

  // Get CPU count on mount
  useEffect(() => {
    wails()?.GetCPUCount().then(n => {
      if (n) { setCpuCount(n); setThreads(Math.min(n * 2, 16)) }
    }).catch(() => {})
  }, [])

  // Register Wails event listeners
  useEffect(() => {
    const runtime = rt()
    if (!runtime) return

    runtime.EventsOn('log', (entry: LogEntry) => {
      setLogs(prev => [...prev, { ...entry, _ts: Date.now() }])
    })

    runtime.EventsOn('progress', (prog: ProgressEvent) => {
      // Only update progress if it carries real data (total > 0)
      // This prevents the blank-emit from Go's defer wiping the UI
      if (prog.total > 0 || prog.processed > 0) {
        setProgress(prog)
      }
      // Track running state from the event stream too
      if (prog.running === false && prog.total > 0) {
        setRunning(false)
        setDone(true)
      }
    })

    return () => {
      runtime.EventsOff('log')
      runtime.EventsOff('progress')
    }
  }, [])

  // Auto-scroll log
  useEffect(() => {
    if (autoScroll && logRef.current) {
      logRef.current.scrollTop = logRef.current.scrollHeight
    }
  }, [logs, autoScroll])

  const handleStart = useCallback(async () => {
    if (!path || running) return
    setLogs([])
    setSummary(null)
    setDone(false)
    setRunning(true)
    setProgress({ ...EMPTY_PROGRESS, running: true })

    try {
      const result = await wails()?.StartOrganize(path, dryRun, threads, excludeExts)
      if (result) {
        setSummary(result)
        // Build progress from summary so stats stay visible after completion
        setProgress(prev => ({
          ...prev,
          total: result.total,
          moved: result.moved,
          skipped: result.skipped,
          errors: result.errors,
          processed: result.total,
          percentDone: 100,
          running: false,
        }))
      }
    } catch (e) {
      console.error(e)
    } finally {
      setRunning(false)
      setDone(true)
    }
  }, [path, dryRun, threads, excludeExts, running])

  const handleCancel = useCallback(async () => {
    await wails()?.Cancel()
    setRunning(false)
  }, [])

  const handleClear = useCallback(() => {
    setLogs([])
    setSummary(null)
    setDone(false)
    setProgress(EMPTY_PROGRESS)
  }, [])

  const catEntries = Object.entries(progress.categoryCounts).sort((a, b) => b[1] - a[1])
  const totalMoved = Object.values(progress.categoryCounts).reduce((a, b) => a + b, 0)

  return (
    <div className="flex flex-col h-screen bg-surface-0 grid-bg overflow-hidden">
      <TitleBar />

      <div className="flex flex-1 overflow-hidden">
        {/* ── LEFT SIDEBAR ──────────────────────────────── */}
        <aside className="w-72 shrink-0 flex flex-col gap-4 p-4 border-r border-white/5 overflow-y-auto">

          <DirectoryPicker path={path} onSelect={setPath} />

          <div className="h-px bg-white/5" />

          <ThreadSlider value={threads} onChange={setThreads} max={cpuCount * 2} />

          <ToggleSwitch checked={dryRun} onChange={setDryRun} label="Dry Run (simulate only)" />

          <div className="space-y-1.5">
            <label className="text-xs font-mono text-slate-500 uppercase tracking-widest">Exclude Extensions</label>
            <input
              type="text"
              value={excludeExts}
              onChange={e => setExcludeExts(e.target.value)}
              placeholder=".tmp, .log, .bak"
              className="no-drag w-full px-3 py-2 rounded-lg bg-surface-2 border border-slate-700 text-xs font-mono text-slate-300
                placeholder-slate-700 focus:outline-none focus:border-cyan-500/50 focus:bg-cyan-500/5 transition-all"
            />
          </div>

          <div className="h-px bg-white/5" />

          {/* Action buttons */}
          <div className="space-y-2">
            {!running ? (
              <button
                onClick={handleStart}
                disabled={!path}
                className={`no-drag w-full py-2.5 rounded-lg font-display tracking-widest text-sm transition-all
                  ${path
                    ? dryRun
                      ? 'bg-amber-500/20 border border-amber-500/50 text-amber-400 hover:bg-amber-500/30'
                      : 'bg-cyan-500/20 border border-cyan-500/50 text-cyan-400 hover:bg-cyan-500/30 glow-cyan'
                    : 'bg-surface-2 border border-white/5 text-slate-700 cursor-not-allowed'
                  }`}
              >
                {dryRun ? '⚡ DRY RUN' : '▶ ORGANIZE FILES'}
              </button>
            ) : (
              <button
                onClick={handleCancel}
                className="no-drag w-full py-2.5 rounded-lg font-display tracking-widest text-sm
                  bg-red-500/20 border border-red-500/50 text-red-400 hover:bg-red-500/30 transition-all"
              >
                ■ CANCEL
              </button>
            )}

            <button
              onClick={handleClear}
              className="no-drag w-full py-2 rounded-lg text-xs font-mono text-slate-600 hover:text-slate-400
                border border-white/5 hover:border-white/10 transition-all"
            >
              Clear
            </button>
          </div>

          <div className="h-px bg-white/5" />

          {/* Category breakdown bars */}
          {catEntries.length > 0 && (
            <div className="space-y-2.5">
              <p className="text-xs font-mono text-slate-600 uppercase tracking-widest">By Category</p>
              {catEntries.map(([cat, count]) => (
                <CategoryBar key={cat} name={cat} count={count} total={totalMoved}
                  color={CATEGORY_COLORS[cat] || '#64748b'} />
              ))}
            </div>
          )}
        </aside>

        {/* ── MAIN CONTENT ───────────────────────────────── */}
        <main className="flex-1 flex flex-col overflow-hidden">

          {/* Progress header */}
          <div className="p-4 border-b border-white/5 space-y-3 shrink-0">
            <ProgressBar percent={progress.percentDone} running={running} done={done && !running} />

            <div className="flex items-center gap-3 flex-wrap">
              <StatPill label="Total"   value={progress.total}   color="#94a3b8" />
              <StatPill label="Moved"   value={progress.moved}   color="#00e676" />
              <StatPill label="Skipped" value={progress.skipped} color="#ffb300" />
              <StatPill label="Errors"  value={progress.errors}  color="#ff5252" />

              {running && progress.currentFile && (
                <div className="ml-auto flex items-center gap-2 min-w-0">
                  <div className="w-1.5 h-1.5 rounded-full bg-cyan-400 animate-pulse-dot shrink-0" />
                  <span className="text-xs font-mono text-slate-500 truncate max-w-xs">
                    {progress.currentFile}
                  </span>
                </div>
              )}

              {done && !running && summary && (
                <div className="ml-auto text-xs font-mono text-green-400 shrink-0">
                  ✓ Done in {summary.elapsedSeconds.toFixed(2)}s
                  {summary.dryRun && <span className="ml-2 text-amber-400">[DRY RUN]</span>}
                </div>
              )}
            </div>
          </div>

          {/* Log feed */}
          <div className="flex-1 relative overflow-hidden">
            <div className="absolute inset-0 flex flex-col">
              {/* Log toolbar */}
              <div className="flex items-center justify-between px-4 py-2 border-b border-white/5 shrink-0">
                <div className="flex items-center gap-2">
                  <span className="text-xs font-mono text-slate-600 uppercase tracking-widest">Live Log</span>
                  {running && <span className="w-1.5 h-1.5 rounded-full bg-cyan-400 animate-pulse-dot" />}
                  <span className="text-xs font-mono text-slate-700">{logs.length} lines</span>
                </div>
                <div className="flex items-center gap-3 no-drag">
                  <button
                    onClick={() => setAutoScroll(v => !v)}
                    className={`text-xs font-mono transition-colors ${autoScroll ? 'text-cyan-500' : 'text-slate-600 hover:text-slate-400'}`}
                  >
                    {autoScroll ? '⇣ Auto-scroll ON' : '⇣ Auto-scroll OFF'}
                  </button>
                  <button
                    onClick={() => setLogs([])}
                    className="text-xs font-mono text-slate-700 hover:text-slate-400 transition-colors"
                  >
                    Clear
                  </button>
                </div>
              </div>

              {/* Log entries */}
              <div
                ref={logRef}
                onScroll={e => {
                  const el = e.currentTarget
                  const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 40
                  setAutoScroll(atBottom)
                }}
                className="flex-1 overflow-y-auto px-4 py-2 space-y-0.5"
              >
                {logs.length === 0 ? (
                  <div className="flex flex-col items-center justify-center h-full text-center gap-3 opacity-30">
                    <svg className="w-12 h-12 text-slate-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1}
                        d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <span className="text-xs font-mono text-slate-600">
                      Select a folder and click Organize to begin
                    </span>
                  </div>
                ) : (
                  logs.map((entry, i) => <LogLine key={i} entry={entry} />)
                )}
              </div>
            </div>
          </div>

          {/* Summary footer — only shown after completion */}
          {summary && !running && (
            <div className="border-t border-white/5 px-4 py-3 bg-surface-1 shrink-0 animate-fade-in">
              <div className="flex items-center gap-4 flex-wrap">
                <span className="font-display tracking-widest text-sm text-green-400">SESSION COMPLETE</span>
                <span className="text-xs font-mono text-slate-600">
                  Moved <span className="text-green-400">{summary.moved}</span> ·
                  Skipped <span className="text-amber-400"> {summary.skipped}</span> ·
                  Errors <span className="text-red-400"> {summary.errors}</span> ·
                  Time <span className="text-cyan-400"> {summary.elapsedSeconds.toFixed(2)}s</span>
                </span>
                {summary.logPath && (
                  <span className="ml-auto text-xs font-mono text-slate-700 truncate max-w-xs" title={summary.logPath}>
                    📄 {summary.logPath}
                  </span>
                )}
              </div>
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
