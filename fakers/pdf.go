package fakers

import (
	"github.com/johnfercher/maroto"
	"os"

	"github.com/adrian-gheorghe/mediafaker/extpoints"
	log "github.com/sirupsen/logrus"
)

func init() {
	faker := new(Pdf)
	faker.Name = "PDF"
	faker.Extensions = []string{"pdf"}
	extpoints.RegisterExtension(faker, "pdf")

}

// Pdf structure
type Pdf struct {
	Name       string
	Extensions []string
}

// Init Method
func (faker *Pdf) init() {
	log.Info("Media Faker Registered: ", faker.Name)
}

// GetExtensions Method
func (faker *Pdf) GetExtensions() []string {
	return faker.Extensions
}

// Fake Command
func (faker *Pdf) Fake(sourcePath string, destinationPath string, sourceInfo os.FileInfo) error {
	log.Info(sourcePath)

	m := maroto.NewMaroto(maroto.Portrait, maroto.A4)
	m.Row(20, func() {
		m.Col(func() {
			m.Text("MediaFaker", &maroto.TextProp{
				Size:  18,
				Align: maroto.Center,
				Top:   17,
			})
		})
	})

	m.Line(1.0)

	err := m.OutputFileAndClose(destinationPath)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// FakeTreeFile Command
func (faker *Pdf) FakeTreeFile(item extpoints.TreeFile, destinationPath string) error {
	log.Info(item.Path, " - faked with PDF")

	m := maroto.NewMaroto(maroto.Portrait, maroto.A4)
	m.Row(20, func() {
		m.Col(func() {
			m.Text("MediaFaker", &maroto.TextProp{
				Size:  18,
				Align: maroto.Center,
				Top:   17,
			})
		})
	})
	m.Line(1.0)

	err := m.OutputFileAndClose(destinationPath)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
