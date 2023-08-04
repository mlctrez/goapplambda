package compo

import (
	"encoding/json"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"net/http"
	"sort"
	"strings"
)

type Headers struct {
	app.Compo
	headerMap map[string]string
}

func (r *Headers) Render() app.UI {

	var rows []app.UI
	if r.headerMap != nil {
		rows = append(rows, app.Tr().Body(
			app.Th().Attr("align", "left").Text("Header"),
			app.Th().Attr("align", "left").Text("Value")),
		)
		var keys []string
		for k := range r.headerMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			rows = append(rows, app.Tr().Body(
				app.Td().Text(k),
				app.Td().Text(r.headerMap[k]),
			))
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
	r.headerMap = m
	r.Update()
}
