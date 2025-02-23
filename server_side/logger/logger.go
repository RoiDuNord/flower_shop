package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

func Init() *slog.Logger {
	logDir, logFile := "logger", "sysLog.log"
	logPath := filepath.Join(logDir, logFile)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("opening log file", "error", err)
		return nil
	}
	defer file.Close()

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(file, opts))

	slog.SetDefault(logger)

	return logger
}

// func getCallerInfo() (string, string, int) {
// 	pc, file, line, ok := runtime.Caller(1)
// 	if !ok {
// 		return "unknown file", "unknown func", 0
// 	}

// 	index := strings.Index(file, "flower_shop")
// 	if index != -1 {
// 		file = file[index:]
// 	} else {
// 		file = "unknown file:" + file
// 	}

// 	fn := runtime.FuncForPC(pc).Name()
// 	fnParts := strings.Split(fn, ".")
// 	if len(fnParts) > 0 {
// 		fn = fnParts[len(fnParts)-1]
// 	} else {
// 		fn = "unknown func:" + fn
// 	}

// 	return file, fn, line
// }
