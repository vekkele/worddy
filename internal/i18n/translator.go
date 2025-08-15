package i18n

import (
	"context"
	"log/slog"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type contextKey string

const i18nCtxKey contextKey = "i18nTranslatorKey"

type M map[string]any

type Translator struct {
	localizer *i18n.Localizer
}

func FromCtx(ctx context.Context) *Translator {
	if service, ok := ctx.Value(i18nCtxKey).(*Translator); ok {
		return service
	}

	return nil
}

func WithTranslator(ctx context.Context, tr *Translator) context.Context {
	return context.WithValue(ctx, i18nCtxKey, tr)
}

func (tr *Translator) T(id string) string {
	return tr.localize(id, nil, nil)
}

func (tr *Translator) Td(id string, data M) string {
	return tr.localize(id, nil, data)
}

func (tr *Translator) N(id string, count int) string {
	return tr.localize(id, &count, nil)
}

func (tr *Translator) Nd(id string, count int, data M) string {
	return tr.localize(id, &count, data)
}

func (tr *Translator) localize(id string, plural *int, data M) string {
	cfg := &i18n.LocalizeConfig{
		MessageID: id,
	}

	if plural != nil {
		if data == nil {
			data = M{"Count": *plural}
		} else if _, ok := data["Count"]; !ok {
			data["Count"] = *plural
		}

		cfg.PluralCount = *plural
	}

	cfg.TemplateData = data

	msg, err := tr.localizer.Localize(cfg)
	if err != nil {
		slog.Warn(err.Error())
		return id
	}

	return msg
}
