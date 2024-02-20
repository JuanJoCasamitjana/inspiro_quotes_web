package quotes

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/", RenderIndex)
	e.GET("/quotes", ListQuotes)
	e.GET("/quotes/:id", FindQuote)
	e.GET("/quotes/:id/images", GetRelatedImages)
	e.GET("/quotes/:id/search-images", SearchRelatedImages)
	e.GET("/generate/image/:id", GeneratePoster)
	e.GET("/genres", FindAllGenres)
	e.GET("/genres/:genre", FindQuotesByGenrePaginated)
}
