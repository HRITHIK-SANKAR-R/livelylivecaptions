package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"strings" 

	"livelylivecaptions/internal/types"
)

// LogLevel defines the severity of a log message.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a single log message.
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
}

// RingBuffer stores log entries in a circular buffer.
type RingBuffer struct {
	entries    []LogEntry
	capacity   int
	head       int // Index of the next write
	count      int // Current number of entries
	mu         sync.RWMutex
	fileLogger *log.Logger // Optional file logger
	minLevel   LogLevel    // Minimum level to log
}

// NewRingBuffer creates a new RingBuffer with a given capacity.
func NewRingBuffer(logConfig types.LogConfig, capacity int) *RingBuffer {
	rb := &RingBuffer{
		entries:  make([]LogEntry, capacity),
		capacity: capacity,
		minLevel: parseLogLevel(logConfig.Level),
	}

	if logConfig.FilePath != "" {
		file, err := os.OpenFile(logConfig.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Failed to open log file %s: %v", logConfig.FilePath, err)
		} else {
			rb.fileLogger = log.New(file, "", log.LstdFlags)
		}
	}
	return rb
}

// Add appends a new log entry to the buffer.
func (rb *RingBuffer) Add(level LogLevel, format string, v ...interface{}) {
	if level < rb.minLevel {
		return // Filter by minLevel
	}

	msg := fmt.Sprintf(format, v...)
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
	}

	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.entries[rb.head] = entry
	rb.head = (rb.head + 1) % rb.capacity
	if rb.count < rb.capacity {
		rb.count++
	}

	if rb.fileLogger != nil {
		rb.fileLogger.Printf("[%s] %s", level.String(), msg)
	}
	// Also print to stderr for immediate feedback unless it's a UI
	if os.Getenv("LIVELY_QUIET_STDOUT") != "1" { // Allow suppressing stdout if needed
		fmt.Fprintf(os.Stderr, "[%s] %s", level.String(), msg)
	}
}

// GetEntries returns all current log entries in chronological order.
func (rb *RingBuffer) GetEntries() []LogEntry {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.count == 0 {
		return nil
	}

	sortedEntries := make([]LogEntry, rb.count)
	if rb.count < rb.capacity {
		copy(sortedEntries, rb.entries[:rb.count])
	} else {
		copy(sortedEntries, rb.entries[rb.head:]) // Copy from head to end
		copy(sortedEntries[rb.capacity-rb.head:], rb.entries[:rb.head]) // Copy from start to head
	}
	return sortedEntries
}

// Global logger instance
var defaultLogger *RingBuffer

// InitGlobalLogger initializes the global logger instance.
// It takes types.LogConfig directly from the main AppConfig.
func InitGlobalLogger(logConfig types.LogConfig) {
	defaultLogger = NewRingBuffer(logConfig, 100) // Default capacity of 100 entries
}

// Helper to parse log level string
func parseLogLevel(levelStr string) LogLevel {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	default:
		return LevelInfo // Default to INFO
	}
}

// Convenience functions for logging
func Debug(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Add(LevelDebug, format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Add(LevelInfo, format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Add(LevelWarn, format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Add(LevelError, format, v...)
	}
}

func Fatal(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Add(LevelFatal, format, v...)
	}
	os.Exit(1) // Fatal logs should exit the application
}
