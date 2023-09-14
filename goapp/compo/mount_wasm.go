package compo

import (
	"bytes"
	"encoding/json"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	fetch "marwan.io/wasm-fetch"
	"strings"
)

func (r *Headers) OnMount(ctx app.Context) {
	target := app.Window().Get("location").Get("href").String()
	target = strings.TrimSuffix(target, "/") + "/api/headers"

	response, err := fetch.Fetch(target, &fetch.Opts{})
	if err != nil {
		app.Log(err)
		return
	}

	m := map[string]string{}
	err = json.NewDecoder(bytes.NewBuffer(response.Body)).Decode(&m)
	if err != nil {
		app.Log(err)
		return
	}

	r.headerMap = m
	r.Update()
}
