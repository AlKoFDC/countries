// +build ignore

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/AlKoFDC/country"
)

var (
	goPkg string
)

func main() {
	/*
	   GOFILE=./countrynames.txt
	   GOPACKAGE=country
	*/
	goFile := os.Getenv("GOFILE")
	goPkg = os.Getenv("GOPACKAGE")
	if goPkg == "" {
		log.Fatal("GOPACKAGE env is empty")
	}
	if goFile == "" {
		log.Fatal("GOFILE env is empty")
	}

	// Parse template for go source
	t, err := template.New("_").Funcs(template.FuncMap{
		"countrySrc":       countrySrc,
		"countryConstName": countryConstName,
	}).Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	// Parse csv data
	f, err := os.Open(goFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	csvr.Comma = ';'
	csvr.Comment = '#'
	csvr.FieldsPerRecord = 27

	var allCountries []country.Country
	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		tr := make([]string, len(record))
		for i, rec := range record {
			tr[i] = strings.TrimSpace(rec)
		}

		country := country.Country{
			Name:         tr[4],
			Alpha2Code:   tr[0],
			Numeric3Code: tr[2],
		}

		allCountries = append(allCountries, country)
	}

	// Prepare template context
	ctx := struct {
		PkgName   string
		Countries []country.Country
	}{
		PkgName:   goPkg,
		Countries: allCountries,
	}

	var buf bytes.Buffer
	// Exec template
	if err := t.Execute(&buf, ctx); err != nil {
		log.Fatal(err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		os.Stderr.Write(buf.Bytes())
		log.Fatalf("fmt: %v", err)
	}

	outf, err := os.Create("countries.gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer outf.Close()

	if _, err = io.Copy(outf, bytes.NewReader(src)); err != nil {
		log.Fatal(err)
	}
}

func countrySrc(country country.Country) string {
	src := strings.Replace(fmt.Sprintf("%#v", country), goPkg+".", "", 1)
	return src
}

func countryConstName(name string) (result string) {
	result = name
	for _, toRemove := range []string{
		" ", ",", "(", ")", "-", "'", ".",
	} {
		result = strings.Replace(result, toRemove, "", -1)
	}
	return
}

const tmpl = `package {{ .PkgName }}

// All countries
var (
{{ range .Countries }}  {{ .Alpha2Code }} = {{ countrySrc . }}
{{ end }}
)

// Numeric codes for all countries as constants.
const (
{{ range .Countries }}  Numeric3{{ countryConstName .Name }} Numeric3Code = "{{ .Numeric3Code }}"
{{ end }}
)

var countryNameMap = map[Name]*Country{
{{ range .Countries }}  "{{ .Name }}" : &{{ .Alpha2Code }},
{{ end }}
}

var iso2LetterMap = map[Alpha2Code]*Country{
{{ range .Countries }}  "{{ .Alpha2Code }}" : &{{ .Alpha2Code }},
{{ end }}
}

var isoNumericMap = map[Numeric3Code]*Country{
{{ range .Countries }}  "{{ .Numeric3Code }}" : &{{ .Alpha2Code }},
{{ end }}
}
`
