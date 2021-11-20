package main

import (
	"encoding/json"
	"fmt"

	"github.com/antchfx/htmlquery"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
)

type Color struct {
	Name string
	Hex  string
}

func main() {
	doc, err := htmlquery.LoadURL("https://www.colorhexa.com/color-names")
	if err != nil {
		log.Error(err)
	}

	var colors []Color

	result := htmlquery.Find(doc, "//table[1]/tbody[1]/tr")

	for _, r := range result {
		trLinks := htmlquery.Find(r, "//a")
		name := trLinks[0].FirstChild.Data
		hex := trLinks[1].FirstChild.Data

		color := Color{
			Name: strcase.ToCamel(name),
			Hex:  hex,
		}

		colors = append(colors, color)
	}

	b, err := json.Marshal(colors)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(string(b))
}
