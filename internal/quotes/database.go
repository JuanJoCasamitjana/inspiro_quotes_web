package quotes

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("quotes.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Quote{}, &SuperUser{}, &ImageForQuote{})

}

func CreateQuote(quote *Quote) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(quote).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetQuote(id uint64) (*Quote, error) {
	var quote Quote
	err := db.First(&quote, id).Error
	return &quote, err
}

func GetQuotesPaginated(page, pageSize int) ([]Quote, error) {
	var quotes []Quote
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&quotes).Error
	return quotes, err
}

func GetQuotesByAuthorPaginated(author string, page, pagesize int) ([]Quote, error) {
	var quotes []Quote
	err := db.Where("author = ?", author).Offset((page - 1) * pagesize).Limit(pagesize).Find(&quotes).Error
	return quotes, err
}

func GetQuotesByGenrePaginated(genre string, page, pagesize int) ([]Quote, error) {
	var quotes []Quote
	err := db.Where("genre = ?", genre).Offset((page - 1) * pagesize).Limit(pagesize).Find(&quotes).Error
	return quotes, err
}

func GetQuotesByLangPaginated(lang string, page, pagesize int) ([]Quote, error) {
	var quotes []Quote
	err := db.Where("lang = ?", lang).Offset((page - 1) * pagesize).Limit(pagesize).Find(&quotes).Error
	return quotes, err
}

func CreateManyQuotes(quotes []*Quote) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&quotes, 200).Error; err != nil {
			return err
		}
		return nil
	})
}

func CreateManyImages(images []*ImageForQuote) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&images, 200).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetImagesByQuoteID(id uint64) ([]ImageForQuote, error) {
	var images []ImageForQuote
	err := db.Where("quote_id = ?", id).Find(&images).Error
	return images, err
}

func GetImageByID(id uint64) (*ImageForQuote, error) {
	var image ImageForQuote
	err := db.Preload("Quote").First(&image, id).Error
	return &image, err
}

func GetAllGenres() ([]string, error) {
	var genres []string
	err := db.Model(&Quote{}).Distinct().Pluck("genre", &genres).Error
	return genres, err
}
