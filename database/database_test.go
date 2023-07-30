package database

import (
	"github.com/enjekt/shannon-engine/models"
	"github.com/enjekt/shannon-engine/pipelines"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	testPan := "5513746525703556"
	db := NewDAO()
	db.Open()
	toEncipher := encipher(testPan, db)

	toDecipher := decipher(toEncipher, db)
	log.Println(toEncipher.ToJSON())
	log.Println(toDecipher.ToJSON())
	assert.Equal(t, testPan, toDecipher.GetPan().String())

}

func decipher(toEncipher models.Palette, db PandaDAO) models.Palette {
	decipherPipeline := pipelines.NewPipeline()
	decipherPipeline.Add(pipelines.DecipherFunc)
	toDecipher := models.NewPalette()
	toDecipher.GetToken().Set(toEncipher.GetToken().String())
	toDecipher.GetPad().Set(toEncipher.GetPad().String())
	db.Get(toDecipher)
	decipherPipeline.Execute(toDecipher)
	return toDecipher
}

func encipher(testPan string, db PandaDAO) models.Palette {
	encipherPipeline := pipelines.NewPipeline()
	encipherPipeline.Add(pipelines.CompactAndStripPanFunc).Add(pipelines.CreatePadFunc).Add(pipelines.EncipherFunc).Add(pipelines.TokenFunc(6, 4))
	toEncipher := models.NewPalette()
	toEncipher.GetPan().Set(testPan)
	toEncipher = encipherPipeline.Execute(toEncipher)

	db.Insert(toEncipher)
	return toEncipher
}
