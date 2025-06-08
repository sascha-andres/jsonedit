package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/google/gops/agent"
	"github.com/sascha-andres/reuse/flag"

	"github.com/sascha-andres/jsonedit"
)

const appName = "JSON_EDIT"

var (
	port      = 8080
	host      = "localhost"
	indent    = "  "
	readOnly  = false
	logLevel  = "info"
	noBrowser = false
)

// init initializes command-line flags for the application,
// setting defaults and descriptions for various configuration options.
func init() {
	flag.SetEnvPrefix(appName)

	flag.IntVar(&port, "port", port, "Port to listen on")
	flag.StringVar(&host, "host", host, "Host to listen on")
	flag.StringVar(&indent, "indent", indent, "Indentation level")
	flag.StringVar(&logLevel, "log-level", logLevel, "Log level (debug, info, warn, error)")
	flag.BoolVar(&readOnly, "read-only", readOnly, "Read-only mode")
	flag.BoolVar(&noBrowser, "no-browaer", noBrowser, "Do not open browser")
}

// main is the entry point of the application, parsing flags and handling any initialization errors during startup.
func main() {
	flag.Parse()

	if err := run(); err != nil {
		panic(err)
	}
}

// run initializes the application with configurations and starts the server,
// returning an error if initialization fails.
func run() error {
	logger := createLogger(logLevel, appName)
	if err := agent.Listen(agent.Options{}); err != nil {
		logger.Error("could not start gops agent", "error", err)
		return err
	}
	a, err := jsonedit.NewApp(
		jsonedit.WithHost(host),
		jsonedit.WithIndent(indent),
		jsonedit.WithPort(port),
		jsonedit.WithReadOnly(readOnly),
		jsonedit.WithLogger(logger),
		jsonedit.WithDebug(logLevel == "debug"),
		jsonedit.WithNoBrowser(noBrowser),
	)
	if err != nil {
		return err
	}
	return a.Run()
}

// CreateLogger initializes and returns a new slog.Logger with the specified log level and project name.
// The log level can be "error", "warn", "info", or "debug".
// Default to "info" if an unknown level is provided.
func createLogger(logLevel, project string) *slog.Logger {
	var handlerOpts *slog.HandlerOptions
	switch strings.ToLower(logLevel) {
	case "warn":
		handlerOpts = &slog.HandlerOptions{Level: slog.LevelWarn}
	case "error":
		handlerOpts = &slog.HandlerOptions{Level: slog.LevelError}
	case "info":
		handlerOpts = &slog.HandlerOptions{Level: slog.LevelInfo}
	case "debug":
		handlerOpts = &slog.HandlerOptions{Level: slog.LevelDebug}
	default:
		handlerOpts = &slog.HandlerOptions{Level: slog.LevelInfo}
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts)).With("project", project)
	slog.SetDefault(logger)
	return logger
}
