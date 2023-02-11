package localization

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"strings"
)

var keys = strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")

type Localizer struct {
	bundle *i18n.Bundle
}

func NewLocalizer(path string, languages ...string) Localizer {
	path = strings.TrimSuffix(path, "/")

	bundle := i18n.NewBundle(language.SimplifiedChinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, lang := range languages {
		bundle.MustLoadMessageFile(fmt.Sprintf("%s/%s.json", path, lang))
	}

	return Localizer{bundle: bundle}
}

func (l Localizer) Get(lang, id string, args []any) string {
	localizer := i18n.NewLocalizer(l.bundle, lang, "zh-CN")

	cfg := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    id,
			Other: id,
		},
		TemplateData: convertI18nTemplateData(args),
	}
	str, err := localizer.Localize(cfg)
	if err != nil {
		return id
	}

	return str
}

// convertI18nTemplateData convert slice of arguments to a map, which is accepted by i18n lib
func convertI18nTemplateData(args []any) map[string]any {
	templateData := make(map[string]any)
	for i, arg := range args {
		templateData[keys[i]] = arg
	}
	return templateData
}
