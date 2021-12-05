package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

func init() {
	DefaultLogger = NewLogger(WithLevel(InfoLevel))
}

type defaultLogger struct {
	sync.RWMutex
	opts Options
}

// Init(opts...) should only overwrite provided options
func (l *defaultLogger) Init(opts ...Option) error {
	for _, o := range opts {
		o(&l.opts)
	}
	return nil
}

func (l *defaultLogger) String() string {
	return "default"
}

func (l *defaultLogger) Fields(fields map[string]interface{}) Logger {
	l.Lock()
	nfields := copyFields(l.opts.Fields)
	l.Unlock()

	for k, v := range fields {
		nfields[k] = v
	}
	dl := &defaultLogger{
		opts: l.opts,
	}
	dl.opts.Fields = nfields
	dl.opts.CallerSkipCount = dl.opts.CallerSkipCount - 1
	return dl
}

func copyFields(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func logCallerfilePath(loggingFilePath string) string {
	// To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	idx := strings.LastIndexByte(loggingFilePath, '/')
	if idx == -1 {
		return loggingFilePath
	}
	idx = strings.LastIndexByte(loggingFilePath[:idx], '/')
	if idx == -1 {
		return loggingFilePath
	}
	return loggingFilePath[idx+1:]
	//parts := strings.Split(loggingFilePath, string(filepath.Separator))
	//return parts[len(parts)-1]
}

func (l *defaultLogger) Log(level Level, v ...interface{}) {
	// TODO decide does we need to write message if log level not used?
	if !l.opts.Level.Enabled(level) {
		return
	}

	l.RLock()
	fields := copyFields(l.opts.Fields)
	l.RUnlock()

	fields["level"] = level.String()

	if _, file, line, ok := runtime.Caller(l.opts.CallerSkipCount); ok {
		fields["file"] = fmt.Sprintf("%s:%d", logCallerfilePath(file), line)
	}

	keys := make([]string, 0, len(fields))
	for k, _ := range fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	metadata := ""

	for _, k := range keys {
		metadata += fmt.Sprintf(" %s=%v", k, fields[k])
	}

	t := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprint(v...)
	fmt.Printf("%s %s %v\n", t, metadata, message)
}

func (l *defaultLogger) Logf(level Level, format string, v ...interface{}) {
	//	 TODO decide does we need to write message if log level not used?
	if level < l.opts.Level {
		return
	}

	l.RLock()
	fields := copyFields(l.opts.Fields)
	l.RUnlock()

	fields["level"] = level.String()

	if _, file, line, ok := runtime.Caller(l.opts.CallerSkipCount); ok {
		fields["file"] = fmt.Sprintf("%s:%d", logCallerfilePath(file), line)
	}

	keys := make([]string, 0, len(fields))
	for k, _ := range fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	metadata := ""

	for _, k := range keys {
		metadata += fmt.Sprintf(" %s=%v", k, fields[k])
	}

	t := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, v...)
	fmt.Printf("%s %s %v\n", t, metadata, message)
}

func (n *defaultLogger) Options() Options {
	// not guard against options Context values
	n.RLock()
	opts := n.opts
	opts.Fields = copyFields(n.opts.Fields)
	n.RUnlock()
	return opts
}

// NewLogger builds a new logger based on options
func NewLogger(opts ...Option) Logger {
	// Default options
	options := Options{
		Level:           InfoLevel,
		Fields:          make(map[string]interface{}),
		Writer:          []io.Writer{os.Stderr},
		CallerSkipCount: 2,
		Context:         context.Background(),
	}

	l := &defaultLogger{opts: options}
	if err := l.Init(opts...); err != nil {
		l.Log(FatalLevel, err)
	}

	return l
}
