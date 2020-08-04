package log

import (
	"time"

	"github.com/liasece/micserver/log/core"
)

// options of log
type options struct {
	// NoConsole while remove console out put, default false.
	NoConsole bool
	// NoConsoleColor whil disable console output color, default false.
	NoConsoleColor bool
	// FilePaths is the log output file path, default none log file.
	FilePaths []string
	// RecordTimeLayout use for (time.Time).Format(layout string) record time field, default "060102-15:04:05",
	// will not be empty.
	RecordTimeLayout string
	// Level log record level limit, only higer thie level log can be get reach Writer, default SYS.
	Level Level
	// Name is thie logger name, default "".
	Name string
	// Topic is thie logger topic, default "".
	Topic string
	// AsyncWrite while asynchronously output the log record to Write, it may be more performance,
	// but if you exit(e.g. os.Exit(1), main() return) this process, it may be loss some log record,
	// because they didn't have time to Write and flush to file.
	AsyncWrite bool
	// AsyncWriteDuration only effective when AsyncWrite is true, this is time duration of asynchronously
	// check log output to write default 100ms.
	AsyncWriteDuration time.Duration
	// RedirectError duplicate stderr to log file, it will be call syscall.Dup2 in linux or syscall.DuplicateHandle
	// in windows, default false.
	RedirectError bool
	// RotateTimeLayout use of (time.Time).Format(layout string) to check if a roteta file is required.
	// default "", will disable rotate. Highest accuracy is minutes.
	RotateTimeLayout string
	// Development use of DPanicLevel panic, default false
	Development bool
	// AddCaller, default false
	AddCaller bool
	// AddStack, default nil
	AddStack core.LevelEnabler
	// CallerSkip, default 2
	CallerSkip int
	// ErrorOutput, default nil
	ErrorOutput core.WriteSyncer
}

var _defaultoptions = options{
	Level:              SysLevel,
	RecordTimeLayout:   "060102-15:04:05",
	AsyncWriteDuration: time.Millisecond * 100,
	CallerSkip:         2,
}

// TOptions type of options
type TOptions struct {
	opts []optionFunc
}

func (o *TOptions) apply(log *Logger) {
	for _, v := range o.opts {
		v.apply(log)
	}
}

// Options func
func Options() *TOptions {
	return &TOptions{}
}

// An Option configures a Logger.
type Option interface {
	apply(*Logger)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

// WrapCore wraps or replaces the Logger's underlying core.Core.
func (o *TOptions) WrapCore(f func(core.Core) core.Core) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.core = f(log.core)
	}))
	return o
}

// Hooks registers functions which will be called each time the Logger writes
// out an Entry. Repeated use of Hooks is additive.
//
// Hooks are useful for simple side effects, like capturing metrics for the
// number of emitted logs. More complex side effects, including anything that
// requires access to the Entry's structured fields, should be implemented as
// a core.Core instead. See core.RegisterHooks for details.
func (o *TOptions) Hooks(hooks ...func(core.Entry) error) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.core = core.RegisterHooks(log.core, hooks...)
	}))
	return o
}

// Fields adds fields to the Logger.
func (o *TOptions) Fields(fs ...Field) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.core = log.core.With(fs)
	}))
	return o
}

// FilePaths add config to the Logger.
func (o *TOptions) FilePaths(path ...string) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.FilePaths = append(log.options.FilePaths, path...)
	}))
	return o
}

// RotateTimeLayout add config to the Logger.
func (o *TOptions) RotateTimeLayout(value string) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.RotateTimeLayout = value
	}))
	return o
}

// NoConsole add config to the Logger.
func (o *TOptions) NoConsole(value bool) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.NoConsole = value
	}))
	return o
}

// NoConsoleColor add config to the Logger.
func (o *TOptions) NoConsoleColor(value bool) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.NoConsoleColor = value
	}))
	return o
}

// AsyncWrite add config to the Logger.
func (o *TOptions) AsyncWrite(value bool) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.AsyncWrite = value
	}))
	return o
}

// Level add config to the Logger.
func (o *TOptions) Level(value Level) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.Level = value
	}))
	return o
}

// Topic add config to the Logger.
func (o *TOptions) Topic(value string) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.Topic = value
	}))
	return o
}

// RecordTimeLayout add config to the Logger.
func (o *TOptions) RecordTimeLayout(value string) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.RecordTimeLayout = value
	}))
	return o
}

// Name add config to the Logger.
func (o *TOptions) Name(value string) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.Name = value
	}))
	return o
}

// AsyncWriteDuration add config to the Logger.
func (o *TOptions) AsyncWriteDuration(value time.Duration) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.AsyncWriteDuration = value
	}))
	return o
}

// RedirectError add config to the Logger.
func (o *TOptions) RedirectError(value bool) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.RedirectError = value
	}))
	return o
}

// AddStack add config to the Logger.
func (o *TOptions) AddStack(value core.LevelEnabler) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.AddStack = value
	}))
	return o
}

// ErrorOutput add config to the Logger.
func (o *TOptions) ErrorOutput(value core.WriteSyncer) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.ErrorOutput = value
	}))
	return o
}

// Development puts the logger in development mode, which makes DPanic-level
// logs panic instead of simply logging an error.
func (o *TOptions) Development() *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.Development = true
	}))
	return o
}

// AddCaller configures the Logger to annotate each message with the filename
// and line number of log's caller.  See also WithCaller.
func (o *TOptions) AddCaller() *TOptions {
	return o.WithCaller(true)
}

// WithCaller configures the Logger to annotate each message with the filename
// and line number of log's caller, or not, depending on the value of enabled.
// This is a generalized form of AddCaller.
func (o *TOptions) WithCaller(enabled bool) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.AddCaller = enabled
	}))
	return o
}

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option). When building wrappers around the
// Logger and SugaredLogger, supplying this Option prevents log from always
// reporting the wrapper code as the caller.
func (o *TOptions) AddCallerSkip(skip int) *TOptions {
	o.opts = append(o.opts, optionFunc(func(log *Logger) {
		log.options.CallerSkip += skip
	}))
	return o
}
