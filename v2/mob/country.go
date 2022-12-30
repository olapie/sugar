package mob

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed resources/countries.json
var countriesJSONString string

type Country struct {
	CallingCode int64  `json:"calling_code"`
	Flag        string `json:"flag"`
	Name        string `json:"name"`
}

func GetCountries() []*Country {
	loadCountries()
	return countries
}

func GetCountryByCallingCode(code int64) *Country {
	loadCountries()
	return codeToCountry[code]
}

var countries []*Country
var codeToCountry map[int64]*Country
var loadCountriesOnce sync.Once

func loadCountries() {
	loadCountriesOnce.Do(func() {
		if err := json.Unmarshal([]byte(countriesJSONString), &countries); err != nil {
			fmt.Println(err)
			countries = []*Country{}
			return
		}

		codeToCountry = make(map[int64]*Country, len(countries))
		for _, c := range countries {
			codeToCountry[c.CallingCode] = c
		}
	})
}
