// Package countries defines country names and codes according to ISO 3166.
// The source is http://opengeocode.org/download/countrynames.txt.
//
// From countrynames.txt header:
//
//      Country Codes to Country Names
//      Created by OpenGeoCode.Org, Submitted into the Public Domain Jan 26, 2014 (version 7)
//
//      Abbreviations:
//      ISO : International Standards Organization
//      BGN : U.S. Board on Geographic Names
//      UNGEGN : United Nations Group of Experts on Geographic Names
//      PCGN : UK Permanent Committee on Geographic Names
//      FAO  : United Nations Food & Agriculture Organization
//      FFO  : German Federal Foreign Office
//
//      Metadata (one entry per line)
//      ISO 3166-1 alpha-2, ISO 3166-1 alpha-3; ISO 3166-1 numeric;
//      ISO 3166-1 English short name (Gazetteer order); ISO 3166-1 English short name (proper reading order); ISO 3166-1 English romanized short name (Gazetteer order); ISO 3166-1 English romanized short name (proper reading oorder);
//      ISO 3166-1 French short name (Gazetteer order); ISO 3166-1 French short name (proper reading order); ISO 3166-1 Spanish short name (Gazetteer order);
//      UNGEGN English formal name; UNGEGN French formal name; UNGEGN Spanish formal name; UNGEGN Russian short and formal name; UNGEGN local short name; UNGEGN local formal name;
//      BGN English short name (Gazetteer order); BGN English short name (proper reading order); BGN English long name; BGN local short name; BGN local long name
//      PCGN English short name (Gazetteer order); PCGN English short name (proper reading order); PCGN English long name; FAO Italian long name; FFO German short name
//
//      NOTES:
//      UNGEGN and BGN local names: when there is more than one language local name, each local name is followed by the 639-1 alpha-2 language code within paranthesis (xx) and separated by a slash (/).
//      Ex. Canada(en)/le Canada(fr)
//
// This is heavily adjusted from the original package at github.com/vincent-petithory/countries for my needs.
package country

import (
	"strconv"
	"strings"
)

//go:generate go run parser.go

type (
	Name         string
	Alpha2Code   string
	Numeric3Code string
)

// Country holds fields for a country as defined by ISO 3166.
type Country struct {
	Name         string
	Alpha2Code   string
	Numeric3Code string
}

// NameToNum converts country name to a country numeric code.
func NameToNum(countryName string) (countryNumCode Numeric3Code, ok bool) {
	allCodes, ok := countryNameMap[Name(strings.ToLower(countryName))]
	if ok {
		countryNumCode = Numeric3Code(allCodes.Numeric3Code)
	}
	return
}

// ISOToNum converts country ISO code to a country numeric code.
func ISOToNum(countryIsoCode Alpha2Code) (countryNumCode Numeric3Code, ok bool) {
	allCodes, ok := iso2LetterMap[countryIsoCode.ToLower()]
	if ok {
		countryNumCode = Numeric3Code(allCodes.Numeric3Code)
	}
	return
}

// NumToISO converts a country numeric code into a country 2 letter ISO code.
func NumToISO(countryNumCode Numeric3Code) (countryISOCode Alpha2Code, ok bool) {
	allCodes, ok := isoNumericMap[Numeric3Code(countryNumCode)]
	if ok {
		countryISOCode = Alpha2Code(strings.ToLower(allCodes.Alpha2Code))
	}
	return
}

// CheckNum validates the country numeric code.
func CheckNum(countryNum Numeric3Code) (countryNumCode Numeric3Code, ok bool) {
	allCodes, ok := isoNumericMap[countryNum]
	if ok {
		countryNumCode = Numeric3Code(allCodes.Numeric3Code)
	}
	return
}

// ToNumeric3 translates the country to the country ISO 3166-1 numeric code.
// It only supports the english country names and ISO 3166-1 alpha-2 codes.
// If country is of length < 2, it will return an empty string and false.
// If country is of length 2, it will interpret it as ISO 2 letter code.
// If country is of length 3, it will assume, that it is already the ISO number, and just check, if it's known.
// If country is of length > 3, it will interpret it as a country name.
// If a number was found, answerOK will be true, false otherwise.
func ToNumeric3(country string) (countryNumCode Numeric3Code, answerOK bool) {
	// It's an ISO 2 letter country code.
	if len(country) == 2 {
		if num, ok := ISOToNum(Alpha2Code(country)); ok {
			countryNumCode = num
			answerOK = true
		}
	}
	// Assuming len(country) == 3 means it's a country ISO number
	if len(country) == 3 {
		if num, ok := CheckNum(Numeric3Code(country)); ok {
			countryNumCode = num
			answerOK = true
		}
	}
	// It's a country name.
	if len(country) > 3 {
		if num, ok := NameToNum(country); ok {
			countryNumCode = num
			answerOK = true
		}
	}
	return
}

// ParseCountries parses country codes from a string list.
func ParseCountries(list string) []Numeric3Code {
	countriesStrings := strings.Split(list, ",")
	countriesMap := make(map[Numeric3Code]bool)
	countriesCodes := make([]Numeric3Code, 0, len(countriesStrings))

	for _, c := range countriesStrings {
		// get rid of spaces if any
		c = strings.Replace(c, " ", "", -1)
		cCode, ok := ToNumeric3(c)
		if ok {
			if !countriesMap[cCode] {
				countriesCodes = append(countriesCodes, cCode)
			}
			countriesMap[cCode] = true
		}
	}
	return countriesCodes
}

// IsValid validates the country numeric 3 code.
func (c Numeric3Code) IsValid() bool {
	if len(c) != 3 {
		return false
	}
	_, err := strconv.Atoi(string(c))
	_, ok := CheckNum(c)
	return err == nil && ok
}

// ToUpper converts an Alpha2Code to its upper case representation
func (a Alpha2Code) ToUpper() Alpha2Code {
	return Alpha2Code(strings.ToUpper(string(a)))
}

// ToLower converts an Alpha2Code to its lower case representation
func (a Alpha2Code) ToLower() Alpha2Code {
	return Alpha2Code(strings.ToLower(string(a)))
}
