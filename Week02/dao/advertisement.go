package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// get db instance
func getDb() *gorm.DB {
	// todo
	panic("no database to connect")
	return nil
}

type Ad struct {
	Id   int
	Name string
	Link string
	Desc string
}

func GetAd(id int) (Ad, error) {
	//return Ad{}, ErrNoRows
	db := getDb()
	ad := Ad{}
	db = db.Select("id = ?", id).First(&ad)

	if err := db.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ad, ErrNoRows
		}
		return ad, errors.Wrap(db.Error, "GetAd Err")
	}
	return ad, nil
}




