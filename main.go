// Main package.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"ghcr.io/srl-labs/ndk-sshx/sshx"
	syslog "github.com/RackSec/srslog"

	"github.com/rs/zerolog"
	"github.com/srl-labs/bond"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	version = "0.0.0"
	commit  = ""
)

// Main entry point for the application.
func main() {
	versionFlag := flag.Bool("version", false, "print the version and exit")

	flag.Parse()

	if *versionFlag {
		fmt.Println(version + "-" + commit)
		os.Exit(0)
	}

	logger := setupLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []bond.Option{
		bond.WithLogger(&logger),
		bond.WithContext(ctx, cancel),
		bond.WithAppRootPath(sshx.AppRoot),
	}

	agent, errs := bond.NewAgent(sshx.AppName, opts...)
	for _, err := range errs {
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create agent")
		}
	}

	err := agent.Start()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start agent")
	}

	app := sshx.New(&logger, agent)
	app.Start(ctx)
}

// setupLogger creates a logger instance.
func setupLogger() zerolog.Logger {
	var writers []io.Writer

	// the lab creates an empty file to indicate
	// that we run in dev mode. If file exists, we
	// log to console as well.
	_, err := os.Stat("/tmp/.ndk-dev-mode")
	if err == nil {
		const logTimeFormat = "2006-01-02 15:04:05 MST"

		consoleLogger := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: logTimeFormat,
			NoColor:    true,
		}

		writers = append(writers, consoleLogger)
	}

	const logFile = "/var/log/sshx/sshx.log"

	// A lumberjack logger with rotation settings.
	fileLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    2, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	}

	var zsyslog zerolog.SyslogWriter
	zsyslog, err = syslog.Dial("", "", syslog.LOG_INFO|syslog.LOG_LOCAL7, "ndk-sshx-go")
	if err != nil {
		panic(err)
	}

	writers = append(writers, fileLogger, zsyslog)

	mw := io.MultiWriter(writers...)

	return zerolog.New(mw).With().Caller().Timestamp().Logger()
}
