package country

import "testing"

//func CountryNameToNum(countryName string) (countryNumCode string, ok bool) {
type CountryNameToNumTestIO struct {
	inName string
	outNum Numeric3Code
	outOK  bool
}

var CountryNameToNumTests = [...]CountryNameToNumTestIO{
	{"", "", false},
	{"deutschland", "", false},
	{"germany", "276", true},
	{"Germany", "276", true},
	{"GERMANY", "276", true},
	{"gerMAny", "276", true},
}

func TestIsValid(t *testing.T) {
	for _, test := range []struct {
		desc   string
		code   Numeric3Code
		result bool
	}{
		{"perfect", "616", true},
		{"too short", "27", false},
		{"too long", "2766", false},
		{"contains other runes", "6l6", false},
	} {
		got := test.code.IsValid()
		if got != test.result {
			t.Fatalf("Expected to get %t for '%s' country code %s validation but got %t", test.result, test.desc, test.code, got)
		}
	}

}

func TestCountryNameToNum(t *testing.T) {
	for _, testIO := range CountryNameToNumTests {
		resultNum, resultOK := NameToNum(testIO.inName)
		if resultNum != testIO.outNum || resultOK != testIO.outOK {
			t.Errorf("For %s\nexpected (%s, %s)\nbut got  (%s, %s)",
				testIO.inName, testIO.outNum, testIO.outOK, resultNum, resultOK)
		}
	}
}

//func CountryISOToNum(countryIsoCode string) (countryNumCode string, ok bool) {
type CountryISOToNumTestIO struct {
	inISO  Alpha2Code
	outNum Numeric3Code
	outOK  bool
}

var CountryISOToNumTests = [...]CountryISOToNumTestIO{
	{"", "", false},
	{"xx", "", false},
	{"de", "276", true},
	{"De", "276", true},
	{"DE", "276", true},
	{"dE", "276", true},
}

func TestCountryISOToNum(t *testing.T) {
	for _, testIO := range CountryISOToNumTests {
		resultNum, resultOK := ISOToNum(testIO.inISO)
		if resultNum != testIO.outNum || resultOK != testIO.outOK {
			t.Errorf("For %s\nexpected (%s, %s)\nbut got  (%s, %s)",
				testIO.inISO, testIO.outNum, testIO.outOK, resultNum, resultOK)
		}
	}
}

type countryNumToISOTestIO struct {
	inNum  Numeric3Code
	outISO Alpha2Code
	outOK  bool
}

var CountryNumToISOTests = [...]countryNumToISOTestIO{
	{"", "", false},
	{"000", "", false},
	{Numeric3Germany, "de", true},
	{Numeric3Turkey, "tr", true},
	{Numeric3Poland, "pl", true},
	{Numeric3France, "fr", true},
}

func TestCountryNumToISO(t *testing.T) {
	for _, testIO := range CountryNumToISOTests {
		resultISO, resultOK := NumToISO(testIO.inNum)
		if resultISO != testIO.outISO || resultOK != testIO.outOK {
			t.Errorf("For %s\nexpected (%s, %s)\nbut got  (%s, %s)",
				testIO.inNum, testIO.outISO, testIO.outOK, resultISO, resultOK)
		}
	}
}

//func ToNumeric3(country string) (countryNumCode string, ok bool) {
type ToNumeric3TestIO struct {
	inCountry string
	outNum    Numeric3Code
	outOK     bool
}

var ToNumeric3Tests = [...]ToNumeric3TestIO{
	{"", "", false},
	{"a", "", false},
	{"xx", "", false},
	{"de", "276", true},
	{"De", "276", true},
	{"DE", "276", true},
	{"dE", "276", true},
	{"ddr", "", false},
	{"123", "", false},
	{"276", "276", true},
	{"deutschland", "", false},
	{"germany", "276", true},
	{"Germany", "276", true},
	{"GERMANY", "276", true},
	{"gerMAny", "276", true},
}

func TestShouldConvertToNumeric3(t *testing.T) {
	for _, testIO := range ToNumeric3Tests {
		resultNum, resultOK := ToNumeric3(testIO.inCountry)
		if resultNum != testIO.outNum || resultOK != testIO.outOK {
			t.Errorf("For %s\nexpected (%s, %t)\nbut got  (%s, %t)",
				testIO.inCountry, testIO.outNum, testIO.outOK, resultNum, resultOK)
		}
	}
}

func TestParseCountriesList(t *testing.T) {
	const list = "DE,TR,PL,CN"
	expected := []Numeric3Code{"276", "792", "616", "156"}
	array := ParseCountries(list)

	if array == nil || len(array) != len(expected) {
		t.Errorf("For %s\nexpected %v\nbut got %v",
			list, expected, array)

	}

	for i := range array {
		if expected[i] != array[i] {
			t.Errorf("For %s\nexpected %v\nbut got %v",
				list, expected, array)
		}
	}
}
