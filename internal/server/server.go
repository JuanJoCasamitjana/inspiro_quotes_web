package server

import (
	"html/template"
	"inspiro_quotes_web/internal/quotes"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupAndRun(port string) {
	env_port := os.Getenv("PORT")
	if env_port != "" {
		port = env_port
	}
	port = ":" + port
	e := echo.New()
	quotes.InitDB()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{method:"${method}", uri:"${uri}", status:"${status}"}
`,
	}))
	e.Renderer = NewTemplates()
	e.Static("/static", "web/static")
	quotes.SetupRoutes(e)
	e.Logger.Fatal(e.Start("0.0.0.0" + port))
}

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("./web/templates/*.html")),
	}
}
