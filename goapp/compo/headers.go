package compo

import (
	"encoding/json"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"net/http"
	"strings"
)

type Headers struct {
	app.Compo
	envMap map[string]string
}

func (r *Headers) Render() app.UI {

	var rows []app.UI
	if r.envMap != nil {
		rows = append(rows, app.Tr().Body(
			app.Th().Attr("align", "left").Text("Header"),
			app.Th().Attr("align", "left").Text("Value")),
		)
		for k, v := range r.envMap {
			rows = append(rows, app.Tr().Body(app.Td().Text(k), app.Td().Text(v)))
		}
		return app.Table().Body(rows...)
	}
	return app.Table().Body()
}

func (r *Headers) OnMount(ctx app.Context) {
	target := app.Window().Get("location").Get("href").String()
	target = strings.TrimSuffix(target, "/") + "/api/headers"

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
