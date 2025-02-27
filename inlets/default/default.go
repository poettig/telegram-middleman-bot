package _default

import (
	"context"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/muety/webhook2telegram/config"
	"github.com/muety/webhook2telegram/model"
	"github.com/muety/webhook2telegram/util"
	"net/http"

	"github.com/muety/webhook2telegram/inlets"
)

type DefaultInlet struct{}

const defaultOrigin = "Webhook2Telegram"

func (i *DefaultInlet) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m model.ExtendedMessage
		var err error

		m, err = i.tryParseBody(r)
		if err != nil {
			m, err = i.tryParseQuery(r)
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if len(m.Origin) == 0 {
			m.Origin = defaultOrigin
		}

		m.Text = "*" + util.EscapeMarkdown(m.Origin) + "* wrote:\n\n" + m.Text

		ctx := r.Context()
		ctx = context.WithValue(ctx, config.KeyMessage, &m.DefaultMessage)
		ctx = context.WithValue(ctx, config.KeyParams, &m.Options)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (i *DefaultInlet) tryParseBody(r *http.Request) (m model.ExtendedMessage, err error) {
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&m)
	return m, err
}

func (i *DefaultInlet) tryParseQuery(r *http.Request) (m model.ExtendedMessage, err error) {
	query := r.URL.Query()
	queryParams := make(map[string]string)
	for k := range r.URL.Query() {
		queryParams[k] = query.Get(k)
	}
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &m,
	})
	err = decoder.Decode(queryParams)
	return m, err
}

func New() inlets.Inlet {
	return &DefaultInlet{}
}
