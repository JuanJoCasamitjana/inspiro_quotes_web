# Quote Poster Generator

This is a simple web application written in Go that generates posters for quotes. Users can input their favorite quotes and customize the appearance of the poster.

## Features

- **Generate Posters**: Select any quote from the database and generate a poster using images provided by [Unsplash](https://unsplash.com/).

## Requirements

- Go 1.21.6 or higher
- Dependencies listed in `go.mod`

## Installation

1. Clone this repository:

```bash
git clone https://github.com/JuanJoCasamitjana/inspiro_quotes_web.git
```

2. Navigate to project directory

```bash
cd inspiro_quotes_web
```

3. Initialize the database

```bash
go env -w CGO_ENABLED=1
go run .\cmd\utils\initialize_db.go -action init
go env -w CGO_ENABLED=0
```

4. Run the application

```bash
go env -w CGO_ENABLED=1
go run .\cmd\utils\main.go
go env -w CGO_ENABLED=0
```

The application runs on port ":8080" and requires an Unsplash api key to fetch the images.
The api keys are accesed through environment variables, but you can set your own ".env" for them.


