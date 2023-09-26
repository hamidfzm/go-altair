package altair

import (
	"embed"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed templates/*
	templatesFS embed.FS
	templates   *template.Template
)

var funcMap = template.FuncMap{
	"getSubscriptionAbsoluteEndpoint": getSubscriptionAbsoluteEndpoint,
	"jsBool":                          strconv.FormatBool,
	"json":                            structToJSONString,
}

const (
	altairTemplate   = "altair.tmpl"
	templatesPattern = "templates/*.tmpl"
)

func init() {
	templates = template.Must(template.New("").Funcs(funcMap).ParseFS(templatesFS, templatesPattern))
}

// Header is a key-value pair of HTTP headers.
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Config is the configuration for the Altair GraphQL IDE.
type Config struct {
	Force              bool
	DefaultWindowTitle string
	Endpoint           string
	Headers            []Header
}

// Handler returns an http.HandlerFunc that serves the Altair GraphQL IDE.
func Handler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.Endpoint = endpointToAbsolute(cfg.Endpoint, r.Host, r.URL.Scheme)
		err := templates.ExecuteTemplate(w, altairTemplate, cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EchoHandler returns an echo.HandlerFunc that serves the Altair GraphQL IDE.
func EchoHandler(cfg *Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		Handler(cfg)(c.Response().Writer, c.Request())
		return nil
	}
}

// getSubscriptionAbsoluteEndpoint returns the subscription absolute endpoint for the given
// endpoint if it is parsable as a URL, or an empty string.
func getSubscriptionAbsoluteEndpoint(endpoint string) string {
	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return ""
	}

	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	default:
		u.Scheme = "ws"
	}

	return u.String()
}

// endpointToAbsolute returns the absolute URL for the given endpoint, using the given host and scheme
func endpointToAbsolute(endpoint string, initialHost string, initialScheme string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		return endpoint
	}

	if u.IsAbs() {
		return endpoint
	}

	if err != nil || u.Scheme != "" {
		return endpoint
	}

	if u.Scheme == "" {
		// If the initial scheme is empty, we default to http.

		if initialScheme == "" {
			initialScheme = "http"
		}
		u.Scheme = initialScheme
	}

	if u.Host == "" {
		u.Host = initialHost
	}

	return u.String()
}

// structToJSONString returns the JSON string representation of the given struct.
func structToJSONString(s any) string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(b)
}
