package xlang

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
)

type keyToString = map[string]string
type globalData struct {
	// e.g. en-US, zh-CN
	Language string
	// e.g. en, zh-Hans
	PrimaryLanguage  string
	LocalizedStrings map[string]keyToString
}

var global globalData

func init() {
	global.LocalizedStrings = make(map[string]keyToString)
}

type UnmarshalFunc = func([]byte, any) error

var unmarshalFuncMap = map[string]UnmarshalFunc{".json": json.Unmarshal}

func SetUnmarshalFunc(extension string, fn UnmarshalFunc) {
	unmarshalFuncMap[extension] = fn
}

func GetUnmarshalFunc(extension string, fn UnmarshalFunc) UnmarshalFunc {
	return unmarshalFuncMap[extension]
}

func SetLang(s string) error {
	t, err := language.Parse(s)
	if err != nil {
		return fmt.Errorf("failed parsing %s, %w", s, err)
	}
	global.Language = t.String()
	switch {
	case t == language.Chinese:
		global.PrimaryLanguage = language.SimplifiedChinese.String()
	case t.Parent() == language.Chinese:
		script, _ := t.Script()
		hans, _ := language.SimplifiedChinese.Script()
		if script.String() == hans.String() {
			global.PrimaryLanguage = language.SimplifiedChinese.String()
		} else {
			global.PrimaryLanguage = language.TraditionalChinese.String()
		}
	case t == language.English || t.Parent() == language.English:
		global.PrimaryLanguage = language.English.String()
	}
	return nil
}

func GetLang() string {
	return global.Language
}

func AddLocalizedStrings(m map[string]map[string]string) {
	for lang, k2s := range m {
		lm := global.LocalizedStrings[lang]
		if len(lm) == 0 {
			global.LocalizedStrings[lang] = k2s
			continue
		}
		for k, s := range k2s {
			lm[k] = s
		}
	}
}

func SetLocalizedStrings(strings map[string]map[string]string) {
	global.LocalizedStrings = strings
}

type ReadDirFileFS interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

func LoadLocalizedStringsFromFS(readFS ReadDirFileFS, dirname string) error {
	entries, err := readFS.ReadDir(dirname)
	if err != nil {
		return fmt.Errorf("fs.ReadDir: %w", err)
	}

	for _, e := range entries {
		filename := filepath.Join(dirname, e.Name())
		data, err := readFS.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("fs.ReadFile: %s, %w", filename, err)
		}
		var m map[string]string
		ext := filepath.Ext(e.Name())
		fn := unmarshalFuncMap[ext]
		if fn == nil {
			return errors.New("unsupported file type")
		}

		if err = fn(data, &m); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		lang := e.Name()[:strings.Index(e.Name(), ".")]
		t, err := language.Parse(lang)
		if err != nil {
			return fmt.Errorf("language.Parse: %s, %w", lang, err)
		}

		if current := global.LocalizedStrings[t.String()]; current == nil {
			global.LocalizedStrings[t.String()] = m
		} else {
			for k, v := range m {
				current[k] = v
			}
		}
	}
	return nil
}

// Localize returns localized string with given key
func Localize(key string) string {
	if m := global.LocalizedStrings[global.Language]; len(m) > 0 {
		tr := m[key]
		if tr != "" {
			return tr
		}
	}

	if m := global.LocalizedStrings[global.PrimaryLanguage]; len(m) > 0 {
		tr := m[key]
		if tr != "" {
			return tr
		}
	}
	return key
}

func IsChinese() bool {
	return global.Language[:2] == "zh"
}

func IsSimplifiedChinese() bool {
	return global.PrimaryLanguage == language.SimplifiedChinese.String()
}

func IsTraditionalChinese() bool {
	return global.PrimaryLanguage == language.SimplifiedChinese.String()
}
