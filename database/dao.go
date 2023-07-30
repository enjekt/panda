package database

import (
	"github.com/enjekt/shannon-engine/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func NewDAO() PandaDAO {
	return new(dao)
}

type PandaDAO interface {
	Open()
	Insert(p models.Palette)
	Get(p models.Palette)
}
type dao struct {
	db *gorm.DB
}

// TODO Pass in in file location...
func (d *dao) Open() {
	d.db, _ = gorm.Open(
		sqlite.Open("panda.db"),
	)
	tp := &TokenPaddedPan{}

	err := d.db.AutoMigrate(&tp)
	if err != nil {
		log.Println(err)
	}
}

func (d *dao) Insert(p models.Palette) {
	tpp := new(TokenPaddedPan)
	tpp.Token = p.GetToken().String()
	tpp.PaddedPan = p.GetPaddedPan().String()
	tx := d.db.Create(tpp)
	log.Println(tx.Error)

}

func (d *dao) Get(p models.Palette) {
	tpp := new(TokenPaddedPan)
	tpp.Token = p.GetToken().String()
	d.db.First(tpp)
	p.GetPaddedPan().Set(tpp.PaddedPan)
}
