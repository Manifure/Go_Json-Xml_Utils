package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Recipes struct {
	XMLName xml.Name `xml:"recipes" json:"-"`
	Cakes   []Cake   `xml:"cake" json:"cake"`
}

type Cake struct {
	Name        string       `xml:"name" json:"name"`
	StoveTime   string       `xml:"stovetime" json:"time"`
	Ingredients []Ingredient `xml:"ingredients>item" json:"ingredients"`
}

type Ingredient struct {
	ItemName  string  `xml:"itemname" json:"ingredient_name"`
	ItemCount float64 `xml:"itemcount" json:"ingredient_count,string"`
	ItemUnit  string  `xml:"itemunit" json:"ingredient_unit"`
}

type DBReader interface {
	Read(data []byte, v interface{}) error
}

func jsonToXml(file []byte) []byte {
	var p Recipes
	err := json.Unmarshal(file, &p)
	if err != nil {
		panic(err)
	}
	p2, err := xml.MarshalIndent(p, "", "    ")
	if err != nil {
		panic(err)
	}
	return p2
}

func xmlToJson(file []byte) []byte {
	var p Recipes
	err := xml.Unmarshal(file, &p)
	if err != nil {
		panic(err)
	}
	p2, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		panic(err)
	}
	return p2
}

func readFile(filename string) []byte {
	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return file
}

func main() {
	filename := flag.String("f", "", "Имя файла для чтения")
	flag.Parse()
	file := readFile(*filename)
	switch filepath.Ext(*filename) {
	case ".xml":
		xmlToJson(file)
		fmt.Println(string(xmlToJson(file)))
	case ".json":
		jsonToXml(file)
		fmt.Println(string(jsonToXml(file)))
	default:
		println("Error: Unsupported file extension")
	}
}
