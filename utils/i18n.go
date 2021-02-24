package utils

import (
	"io/ioutil"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v2"
)

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	directory := "i18n"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".yml" {
			bundle.LoadMessageFile(directory + "/" + f.Name())
		}
	}
}

func Translate(lang string, key string) string {
	return i18n.NewLocalizer(bundle, lang).MustLocalize(&i18n.LocalizeConfig{MessageID: key})
}

func TranslateWithArgs(lang string, key string, args map[string]string) string {
	return i18n.NewLocalizer(bundle, lang).MustLocalize(&i18n.LocalizeConfig{MessageID: key, TemplateData: args})
}
