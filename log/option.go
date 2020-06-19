package log

import "time"

// Options of log
type Options struct {
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
}

var defaultOptions = Options{
	RecordTimeLayout: "060102-15:04:05",
}

// Check options
func (o *Options) Check() {
	if o.RecordTimeLayout == "" {
		o.RecordTimeLayout = "060102-15:04:05"
	}
	if o.AsyncWrite && o.AsyncWriteDuration.Milliseconds() == 0 {
		o.AsyncWriteDuration = time.Millisecond * 100
	}
}
