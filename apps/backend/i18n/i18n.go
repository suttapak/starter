package i18n

import (
	"context"
	"embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	goI8n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/logger"
	"golang.org/x/text/language"
)

//go:embed active.*.json
var LocaleFS embed.FS

func newI18n(conf *config.Config, log logger.AppLogger) (*goI8n.Bundle, error) {
	bundle := goI8n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	_, err := bundle.LoadMessageFileFS(LocaleFS, "active.en.json")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	_, err = bundle.LoadMessageFileFS(LocaleFS, "active.th.json")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return bundle, nil
}

type Local string

const (
	TH Local = "th"
	EN Local = "en"
)

func (i Local) IsValid() bool {
	return i == "th" || i == "en"
}

type (
	I18N interface {
		GetMessage(local Local, id string) string
	}
	i18n struct {
		i18n *goI8n.Bundle
	}
)

// GetMessage implements I18N.
func (i *i18n) GetMessage(local Local, id string) string {
	if !local.IsValid() {
		return "###"
	}
	localizer := goI8n.NewLocalizer(i.i18n, string(local), "th", "en")
	message := localizer.MustLocalize(&goI8n.LocalizeConfig{
		MessageID: id,
	})
	return message
}

func NewI18N(conf *config.Config, log logger.AppLogger) (I18N, error) {
	bundle, err := newI18n(conf, log)
	if err != nil {
		return nil, err
	}
	return &i18n{bundle}, nil
}

func GetLocal(ctx context.Context) (local Local, err error) {
	c, ok := ctx.(*gin.Context)
	if !ok {
		return TH, nil
	}
	l, ok := c.Get("lng")
	if !ok {
		return TH, nil
	}
	local, ok = l.(Local)
	if !ok {
		return TH, nil
	}
	return local, nil
}

func SetLocal(c *gin.Context) {
	local := Local(c.GetHeader("lng"))
	if !local.IsValid() {
		c.Set("lng", TH)
		c.Next()
		return
	}
	c.Set("lng", local)
	c.Next()
}
