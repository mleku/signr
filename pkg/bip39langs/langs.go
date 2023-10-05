package bip39langs

import (
	"github.com/islishude/bip39"
	"sort"
)

var Map = map[string]bip39.Language{
	"english":             bip39.English,
	"czech":               bip39.Czech,
	"chinese-traditional": bip39.ChineseTraditional,
	"chinese-simplified":  bip39.ChineseSimplified,
	"french":              bip39.French,
	"italian":             bip39.Italian,
	"japanese":            bip39.Japanese,
	"korean":              bip39.Korean,
	"portuguese":          bip39.Portuguese,
	"spanish":             bip39.Spanish,
}

func GetList() (langList string) {
	var langs []string
	for i := range Map {
		langs = append(langs, i)
	}
	sort.Strings(langs)
	langList = "[ "
	for i := range langs {
		langList += langs[i]
		if i < len(langs)-1 {
			langList += " | "
		} else {
			langList += " ]"
		}
	}
	return
}
