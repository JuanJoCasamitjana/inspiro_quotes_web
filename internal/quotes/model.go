package quotes

import (
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type Quote struct {
	ID     uint64
	Text   string
	Author string
	Lang   string
	Genre  string
}

type ImageForQuote struct {
	ID       uint64
	QuoteID  uint64
	Quote    Quote
	Author   string
	Source   string `gorm:"default:'unsplash'"`
	ImageURL string
}

type SuperUser struct {
	ID       uint64
	Username string `gorm:"unique"`
	Password string
}

func (s *SuperUser) ValidateAndHashPassword(pass string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return err
	}
	encodedHashedPassword := hex.EncodeToString(hashedPass)
	s.Password = encodedHashedPassword
	return nil
}

func (s *SuperUser) ComparePassword(pass string) bool {
	passBytes, err := hex.DecodeString(s.Password)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(passBytes, []byte(pass))
	return err == nil
}

func (s *SuperUser) SetPasswordFromBytes(pass []byte) error {
	hashedPass, err := bcrypt.GenerateFromPassword(pass, 10)
	if err != nil {
		return err
	}
	encodedHashedPassword := hex.EncodeToString(hashedPass)
	s.Password = encodedHashedPassword
	return nil
}
