package jsonedit

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/doganarif/govisual"
)

type (
	// Define data structures for templates
	EditPageData struct {

		// Content represents the primary JSON content to be displayed or edited on the page.
		Content string

		// Error represents an optional error message to be displayed on the page.
		Error string

		// FormContent contains HTML content used to render an editable JSON form or dynamic interface for the edit page.
		FormContent template.HTML

		// ReadOnly indicates whether the content or form should be displayed in a non-editable mode.
		ReadOnly bool
	}

	// App represents the core application with configuration options such as host, port, JSON indentation, and logging.
	App struct {

		// port specifies the TCP port on which the application server will listen for incoming HTTP requests.
		port int

		// host specifies the hostname or IP address on which the application server will listen for incoming HTTP requests.
		host string

		// indent defines the string used for spacing/indentation when formatting JSON output.
		indent string

		// readOnly indicates whether the application is in read-only mode, restricting modification of data or settings.
		readOnly bool

		// debug indicates whether the application runs in debug mode, enabling detailed request/response logging for development.
		debug bool

		// logger provides structured logging capabilities for the application, supporting various log levels and outputs.
		logger *slog.Logger

		// noBrowser indicates whether the application should open the default browser after startup.
		noBrowser bool
	}

	// AppOption represents a function that configures an App instance and may return an error during the setup process.
	AppOption func(*App) error
)

// NewApp initializes a new App instance with provided AppOption configurations
// and returns the configured App or an error.
func NewApp(opts ...AppOption) (*App, error) {
	app := &App{
		port:      8080,
		host:      "localhost",
		indent:    "  ",
		readOnly:  false,
		logger:    slog.Default(),
		noBrowser: false,
	}
	for _, opt := range opts {
		err := opt(app)
		if err != nil {
			return nil, err
		}
	}
	return app, nil
}

// WithNoBrowser sets the noBrowser field in the App to determine if the app should refrain from opening a browser on startup.
func WithNoBrowser(noBrowser bool) AppOption {
	return func(app *App) error {
		app.noBrowser = noBrowser
		return nil
	}
}

// WithDebug sets the debug mode for the application, enabling detailed logging for requests and responses.
func WithDebug(debug bool) AppOption {
	return func(app *App) error {
		app.debug = debug
		return nil
	}
}

// WithLogger sets the application's logger for structured logging and returns an AppOption for configuration.
func WithLogger(logger *slog.Logger) AppOption {
	return func(app *App) error {
		app.logger = logger
		return nil
	}
}

// WithReadOnly sets the application's read-only mode to the specified value.
func WithReadOnly(readOnly bool) AppOption {
	return func(app *App) error {
		app.readOnly = readOnly
		return nil
	}
}

// WithPort sets the TCP port for the application server to listen on and returns an AppOption for configuration.
func WithPort(port int) AppOption {
	return func(app *App) error {
		app.port = port
		return nil
	}
}

// WithHost sets the hostname or IP address that the application server will use to listen for incoming HTTP requests.
func WithHost(host string) AppOption {
	return func(app *App) error {
		app.host = host
		return nil
	}
}

// WithIndent sets the string used for JSON output indentation in the application configuration.
func WithIndent(indent string) AppOption {
	return func(app *App) error {
		app.indent = indent
		return nil
	}
}

// openBrowser opens the default browser with the specified URL.
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		slog.Error("Failed to open browser", "error", err)
	}
}

// Run starts the HTTP server, register route handlers, and optionally enables debug mode with request/response logging.
func (app *App) Run() error {

	mux := http.NewServeMux()

	// Serve static files from the embedded filesystem
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(GetEmbeddedFileSystem())))

	mux.HandleFunc("/", app.handleIndex)
	mux.HandleFunc("/upload", app.handleUpload)
	mux.HandleFunc("/save", app.handleSave)
	mux.HandleFunc("/edit", app.handleEdit)
	mux.HandleFunc("/new", app.handleNewObject)
	mux.HandleFunc("/new-array", app.handleNewArray)
	mux.HandleFunc("/compare", app.handleCompare)
	mux.HandleFunc("/flatten", app.handleFlatten)
	mux.HandleFunc("/from-schema", app.handleFromSchema)
	mux.HandleFunc("/validate", app.handleValidate)

	app.logger.Info("server starting", "host", app.host, "port", app.port)

	if !app.noBrowser {
		// Open browser after a short delay to ensure server is ready
		go func() {
			time.Sleep(500 * time.Millisecond)
			url := fmt.Sprintf("http://%s:%d", app.host, app.port)
			app.logger.Info("opening browser", "url", url)
			openBrowser(url)
		}()
	}

	if app.debug {
		// Wrap with GoVisual
		handler := govisual.Wrap(
			mux,
			govisual.WithRequestBodyLogging(true),
			govisual.WithResponseBodyLogging(true),
		)

		return http.ListenAndServe(fmt.Sprintf("%s:%d", app.host, app.port), handler)
	}
	return http.ListenAndServe(fmt.Sprintf("%s:%d", app.host, app.port), mux)
}
