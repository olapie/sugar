package xlang

import (
	"embed"
	"testing"
)

//go:embed testdata/*
var testdataFS embed.FS

func TestSetLang(t *testing.T) {
	t.Run("Good", func(t *testing.T) {
		tests := []string{"en", "en-us", "en-US", "zh", "zh-CN", "zh-TW", "zh-HK", "zh-sg", "zh-cn", "zh-hk"}
		for _, test := range tests {
			if err := SetLang(test); err != nil {
				t.Fatal(test, err)
			}
			if GetLang() == "" {
				t.Fatal("no language set")
			}
		}
	})

	t.Run("Bad", func(t *testing.T) {
		SetLang("en")
		tests := []string{"en1", "en-us1", "en-US2", "zh1", "zh-CN2", "zh-TW1", "zh-HK2", "zh-sg1", "zh-cn1", "zh-hk1"}
		for _, test := range tests {
			if err := SetLang(test); err == nil {
				t.Fatal(test)
			}
			if GetLang() != "en" {
				t.Fatal("default language is not en")
			}
		}
	})
}

func TestLoadEmbed(t *testing.T) {
	err := LoadFS(testdataFS, "testdata")
	if err != nil {
		t.Fatal("load embed testdata", err)
	}
	err = SetLang("zh-hans")
	if err != nil {
		t.Fatal("set language", err)
	}
	if tr := Translate("noodle"); tr != "麵" {
		t.Fatal("noodle is incorrect", tr)
	}
	if tr := Translate("noodles"); tr != "noodles" {
		t.Fatal("noodles shouldn't have any translation")
	}
	SetLang("zh-hant")
	if tr := Translate("noodle"); tr != "面" {
		t.Fatal("noodle is incorrect", tr)
	}
}
