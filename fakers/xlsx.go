package fakers

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"os"

	"github.com/adrian-gheorghe/mediafaker/extpoints"
	log "github.com/sirupsen/logrus"
)

func init() {
	faker := new(Xlsx)
	faker.Name = "XLSX"
	faker.Extensions = []string{"xlsx"}
	extpoints.RegisterExtension(faker, "xlsx")

}

// Xlsx structure
type Xlsx struct {
	Name       string
	Extensions []string
}

// Init Method
func (faker *Xlsx) init() {
	log.Info("Media Faker Registered: ", faker.Name)
}

// GetExtensions Method
func (faker *Xlsx) GetExtensions() []string {
	return faker.Extensions
}

// Fake Command
func (faker *Xlsx) Fake(sourcePath string, destinationPath string, sourceInfo os.FileInfo) error {
	log.Info(sourcePath)

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "MediaFaker"
	err = file.Save(destinationPath)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// FakeTreeFile Command
func (faker *Xlsx) FakeTreeFile(item extpoints.TreeFile, destinationPath string) error {
	log.Info(item.Path, " - faked with XLSX")

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "MediaFaker"
	err = file.Save(destinationPath)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
