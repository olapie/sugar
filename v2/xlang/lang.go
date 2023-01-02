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

var langName = "en"
var defaultLangName = "en"
var localizedStrings = map[string]map[string]string{}

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
	langName = t.String()
	switch {
	case t == language.Chinese:
		defaultLangName = language.SimplifiedChinese.String()
	case t.Parent() == language.Chinese:
		script, _ := t.Script()
		hans, _ := language.SimplifiedChinese.Script()
		if script.String() == hans.String() {
			defaultLangName = language.SimplifiedChinese.String()
		} else {
			defaultLangName = language.TraditionalChinese.String()
		}
	case t == language.English || t.Parent() == language.English:
		defaultLangName = language.English.String()
	}
	return nil
}

func GetLang() string {
	return langName
}

func SetLocalizedStrings(strings map[string]map[string]string) {
	localizedStrings = strings
}

type ReadDirFileFS interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

func LoadFS(readFS ReadDirFileFS, dirname string) error {
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
		if current := localizedStrings[t.String()]; current == nil {
			localizedStrings[t.String()] = m
		} else {
			for k, v := range m {
				current[k] = v
			}
		}
	}
	return nil
}

func Translate(s string) string {
	if m := localizedStrings[langName]; len(m) > 0 {
		tr := m[s]
		if tr != "" {
			return tr
		}
	}

	if m := localizedStrings[defaultLangName]; len(m) > 0 {
		tr := m[s]
		if tr != "" {
			return tr
		}
	}
	return s
}
