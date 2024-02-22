package quotes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var ACCESS_KEY string
var SECRET_KEY string

var API_URL string = "https://api.unsplash.com/photos/random"

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	ACCESS_KEY = os.Getenv("ACCESS_KEY")
	SECRET_KEY = os.Getenv("SECRET_KEY")
}

func LoadQuotes(c echo.Context) error {
	return c.String(200, "Quotes will be loaded here")
}

func FindQuote(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Invalid ID")
	}
	quote, err := GetQuote(id)
	if err != nil {
		return c.String(404, "Quote not found")
	}
	return c.Render(200, "quote", quote)
}

func ListQuotes(c echo.Context) error {
	page := c.QueryParam("page")
	pageSize := 10
	if page == "" {
		page = "1"
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		return c.String(400, "Invalid page number")
	}
	if p < 1 {
		return c.String(400, "Invalid page number")
	}
	quotes, err := GetQuotesPaginated(p, pageSize)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	data := map[string]interface{}{
		"quotes": quotes,
		"more":   len(quotes) == pageSize,
		"next":   p + 1,
	}
	return c.Render(200, "quotes", data)
}

func RenderIndex(c echo.Context) error {
	return c.Render(200, "index", nil)
}

func SearchRelatedImages(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Invalid ID")
	}
	quote, err := GetQuote(id)
	if err != nil {
		return c.String(404, "Quote not found")
	}
	urls := requestImages(quote.Text)
	var images []*ImageForQuote
	for _, value := range urls {
		var image ImageForQuote
		image.QuoteID = id
		image.ImageURL = value["imageURL"]
		image.Source = value["source"]
		image.Author = value["author"]
		images = append(images, &image)
	}
	CreateManyImages(images)
	saved_images, err := GetImagesByQuoteID(id)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	var data = make([]map[string]interface{}, len(saved_images))
	for i, img := range saved_images {
		data[i] = map[string]interface{}{
			"imageURL": img.ImageURL,
			"source":   img.Source,
			"author":   img.Author,
			"id":       img.ID,
		}
	}
	return c.Render(200, "images", data)

}

func GetRelatedImages(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Invalid ID")
	}
	images, err := GetImagesByQuoteID(id)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	var data = make([]map[string]interface{}, len(images))
	for i, img := range images {
		data[i] = map[string]interface{}{
			"imageURL": img.ImageURL,
			"source":   img.Source,
			"author":   img.Author,
			"id":       img.ID,
		}
	}
	return c.Render(200, "images", data)
}

func GeneratePoster(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Invalid ID")
	}
	image, err := GetImageByID(id)
	if err != nil {
		return c.String(404, "Image not found")
	}
	body := getImageFromUrl(image.ImageURL)
	if body == nil {
		return c.String(500, "Internal server error")
	}
	generated_poster := editImage(body, image.Quote.Text, image.Quote.Author)
	if generated_poster == nil {
		return c.String(500, "Internal server error")
	}
	encoded_image := base64.StdEncoding.EncodeToString(generated_poster)
	data := map[string]interface{}{
		"image":  encoded_image,
		"quote":  image.Quote.Text,
		"author": image.Quote.Author,
	}
	return c.Render(200, "poster", data)
}

func FindAllGenres(c echo.Context) error {
	genres, err := GetAllGenres()
	if err != nil {
		return c.String(500, "Internal server error")
	}
	return c.Render(200, "genres", genres)
}

func FindQuotesByGenrePaginated(c echo.Context) error {
	genre := c.Param("genre")
	page := c.QueryParam("page")
	pageSize := 10
	if page == "" {
		page = "1"
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		return c.String(400, "Invalid page number")
	}
	if p < 1 {
		return c.String(400, "Invalid page number")
	}
	quotes, err := GetQuotesByGenrePaginated(genre, p, pageSize)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	data := map[string]interface{}{
		"quotes": quotes,
		"more":   len(quotes) == pageSize,
		"next":   p + 1,
		"genre":  genre,
	}
	return c.Render(200, "quotes", data)
}

func requestImages(qt string) []map[string]string {
	query_params := url.Values{}
	query_params.Add("query", qt)
	query_params.Add("orientation", "landscape")
	query_params.Add("topics", "nature, cool tones")
	query_params.Add("count", "6")
	req_url := API_URL + "?" + query_params.Encode()
	req, err := http.NewRequest("GET", req_url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Client-ID "+ACCESS_KEY)
	client := &http.Client{}
	respone, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer respone.Body.Close()
	var data []map[string]interface{}
	if err := json.NewDecoder(respone.Body).Decode(&data); err != nil {
		return nil
	}
	urls := make([]map[string]string, 6)
	for i, v := range data {
		urls[i] = map[string]string{
			"imageURL": v["urls"].(map[string]interface{})["regular"].(string),
			"source":   "unsplash",
			"author":   v["user"].(map[string]interface{})["username"].(string),
		}
	}
	return urls
}

func getImageFromUrl(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	return body
}

func editImage(image []byte, quote, author string) []byte {
	img, err := jpeg.Decode(bytes.NewReader(image))
	if err != nil {
		return nil
	}
	if err != nil {
		return nil
	}
	dc := gg.NewContextForImage(img)
	fontPath := "./web/fonts/" + select_a_random_font()
	err = dc.LoadFontFace(fontPath, 40)
	if err != nil {
		return nil
	}
	width := float64(dc.Width())
	height := float64(dc.Height())
	max_value := 255
	min_value := 230

	randomR := rand.Intn(max_value-min_value+1) + min_value
	randomG := rand.Intn(max_value-min_value+1) + min_value
	randomB := rand.Intn(max_value-min_value+1) + min_value

	dc.SetRGBA255(randomR, randomG, randomB, 70)

	dc.DrawRectangle(0, 0, float64(dc.Width()), float64(dc.Height()))
	dc.Fill()

	maxLineWidth := width * 0.8

	lines := wordWrap(quote, maxLineWidth, dc)

	lineHeight := 40 // Altura estimada de una línea de texto
	startY := (height - float64(len(lines))*float64(lineHeight)) / 2

	dc.SetColor(color.Black)

	for _, line := range lines {
		lineWidth, textHeight := dc.MeasureString(line)

		textX := (width - lineWidth) / 2
		textY := startY

		dc.DrawStringAnchored(line, textX, textY, 0, 0)

		startY += float64(textHeight) + 10
	}
	lineWidth, _ := dc.MeasureString(author)
	textX := (width - lineWidth) / 2
	textY := startY
	dc.DrawStringWrapped(author, textX, textY, 0, 0, maxLineWidth, 1, gg.AlignCenter)

	resImg := dc.Image()

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resImg, nil)
	if err != nil {
		return nil
	}
	image = buf.Bytes()
	return image
}

func wordWrap(text string, maxWidth float64, dc *gg.Context) []string {
	words := strings.Fields(text)
	var lines []string
	var line string
	lineWidth := 0.0

	for _, word := range words {
		wordWidth, _ := dc.MeasureString(word)
		if lineWidth+wordWidth <= maxWidth {
			if line != "" {
				line += " "
			}
			line += word
			lineWidth += wordWidth
		} else {
			lines = append(lines, line)

			line = word
			lineWidth = wordWidth
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func select_a_random_font() string {
	directory := "./web/fonts"
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Println(err)
	}
	max_index := len(files) - 1
	random_index := rand.Intn(max_index)
	selected_font := files[random_index].Name()
	return selected_font
}
