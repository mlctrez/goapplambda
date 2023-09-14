package compo

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"sort"
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
