package txt

import (
	"main/logger"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var once sync.Once
var instance *Txt

type Txt struct {
	bundle *i18n.Bundle
	local  *i18n.Localizer
}

// TODO Mutex
func GetInstance() *Txt {
	once.Do(func() {
		instance = &Txt{}
		instance.bundle = i18n.NewBundle(language.English)

		instance.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

		instance.bundle.ParseMessageFileBytes(LangEnToml.StaticContent, LangEnToml.StaticName)
		instance.bundle.ParseMessageFileBytes(LangRoToml.StaticContent, LangRoToml.StaticName)
	})
	return instance
}

func (l *Txt) SetLanguage(langCode string) *Txt {
	instance.local = i18n.NewLocalizer(instance.bundle, langCode)
	return instance
}

func (l *Txt) GetLabel(key string) string {
	msg, tag, err := l.local.LocalizeWithTag(&i18n.LocalizeConfig{MessageID: key})
	if err != nil {
		logger.Log.Warnf("Could not extract label with key %s language %s because: %v", key, tag, err)
		msg = key
	}

	return msg
}

func GetLabel(key string) string {
	return GetInstance().GetLabel(key)
}
