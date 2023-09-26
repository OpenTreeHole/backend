package data

import (
	"embed"
	_ "embed"
	"os"

	"github.com/goccy/go-json"
)

//go:embed names.json
var NamesFile []byte

//go:embed cedict_ts.u8
var CreditTs []byte

//go:embed keys
var Keys embed.FS

var NamesMapping map[string]string

func init() {
	NamesMappingData, err := os.ReadFile(`data/names_mapping.json`)
	if err != nil {
		return
	}

	err = json.Unmarshal(NamesMappingData, &NamesMapping)
	if err != nil {
		return
	}
}
