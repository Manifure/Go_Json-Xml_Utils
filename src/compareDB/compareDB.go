package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/r3labs/diff/v3"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func jsonUnmarshal(file []byte) Recipes {
	var p Recipes
	err := json.Unmarshal(file, &p)
	if err != nil {
		panic(err)
	}
	return p
}

func xmlUnmarshal(file []byte) Recipes {
	var p Recipes
	err := xml.Unmarshal(file, &p)
	if err != nil {
		panic(err)
	}
	return p
}

func readFile(filename string) []byte {
	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return file
}

func detectExtension(filename *string) Recipes {
	var r Recipes
	file := readFile(*filename)
	switch filepath.Ext(*filename) {
	case ".xml":
		r = xmlUnmarshal(file)
	case ".json":
		r = jsonUnmarshal(file)

	default:
		println("Error: Unsupported file extension")
	}
	return r
}

func pairs(p []string, r *Recipes) string {
	pairs := make([]string, len(p)/2+len(p)%2)
	var a, b int
	for a = len(pairs) - 1; b < len(p)-1; b, a = b+2, a-1 {
		idx, err := strconv.Atoi(p[b+1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if p[b] == "Cakes" {
			pairs[a] = fmt.Sprintf("%s %s", p[b], r.Cakes[idx].Name)
		} else {
			idx1, err := strconv.Atoi(p[1])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			pairs[a] = fmt.Sprintf("%s %s", p[b], r.Cakes[idx1].Ingredients[idx1].ItemName)
		}
	}
	if a == 0 {
		pairs[a] = p[b]
	}
	return strings.Join(pairs, " for ")
}

func CompareRecipes(old *Recipes, new *Recipes) {
	differ, err := diff.NewDiffer(diff.DisableStructValues(), diff.SliceOrdering(false))
	if err != nil {
		panic(err)
	}
	log, err := differ.Diff(old, new)
	if err != nil {
		panic(err)
	}
	for _, change := range log {
		a := change.Path
		if a[0] == "XMLName" {
			continue
		}
		switch change.Type {
		case diff.CREATE:
			fmt.Printf("ADDED %s\n", pairs(a, new))
		case diff.UPDATE:
			fmt.Printf("CHANGED %s - %s instead of %s\n", pairs(a, new), change.To, change.From)
		case diff.DELETE:
			switch n := len(a) - 1; a[n] {
			case "unit":
				a = append(a, change.From.(string))
			case "ingredient":
				a = a[:n]
			}
			fmt.Printf("REMOVED %s\n", pairs(a, old))
		}
	}
}

func main() {
	oldFilename := flag.String("old", "", "Имя файла для чтения")
	newFilename := flag.String("new", "", "Имя файла для чтения")
	flag.Parse()
	old := detectExtension(oldFilename)
	new0 := detectExtension(newFilename)
	CompareRecipes(&old, &new0)
}
