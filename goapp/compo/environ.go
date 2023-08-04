package compo

import (
	"encoding/json"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"net/http"
	"strings"
)

type Environ struct {
	app.Compo
	envMap map[string]string
}

func (r *Environ) Render() app.UI {

	var rows []app.UI
	if r.envMap != nil {
		rows = append(rows, app.Tr().Body(
			app.Th().Attr("align", "left").Text("ENV_VAR"),
			app.Th().Attr("align", "left").Text("VALUE")),
		)
		for k, v := range r.envMap {
			rows = append(rows, app.Tr().Body(app.Td().Text(k), app.Td().Text(v)))
		}
		return app.Table().Body(rows...)
	}
	return app.Table().Body()
}

func (r *Environ) OnMount(ctx app.Context) {
	target := app.Window().Get("location").Get("href").String()
	target = strings.TrimSuffix(target, "/") + "/api/environment"

	resp, err := http.Get(target)
	if err != nil {
		app.Log(err)
		return
	}
	m := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		app.Log(err)
		return
	}
	r.envMap = m
	r.Update()
}
