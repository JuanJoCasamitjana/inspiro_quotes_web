package server

import (
	"html/template"
	"inspiro_quotes_web/internal/quotes"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupAndRun() {
	e := echo.New()
	quotes.InitDB()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{method:"${method}", uri:"${uri}", status:"${status}"}
`,
	}))
	e.Renderer = NewTemplates()
	e.Static("/static", "web/static")
	quotes.SetupRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
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
