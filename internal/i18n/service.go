package i18n

import (
	"embed"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type LocaleService struct {
	bundle *i18n.Bundle
}

const localesPath = "locales"

//go:embed locales
var localeFS embed.FS

func NewService(defaultLang language.Tag) (*LocaleService, error) {
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	files, err := localeFS.ReadDir("locales")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		path := path.Join(localesPath, file.Name())
		if file.IsDir() {
			continue
		}

		buf, err := localeFS.ReadFile(path)
		if err != nil {
			return nil, err
		}

		_, err = bundle.ParseMessageFileBytes(buf, path)
		if err != nil {
			return nil, err
		}
	}

	return &LocaleService{
		bundle: bundle,
	}, nil
}

func (ls *LocaleService) CreateTranslator(locales ...string) *Translator {
	localizer := i18n.NewLocalizer(ls.bundle, locales...)

	return &Translator{
		localizer: localizer,
	}
}
