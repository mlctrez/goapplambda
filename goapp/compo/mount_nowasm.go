//go:build !wasm

package compo

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func (r *Headers) OnMount(ctx app.Context) {
	// does nothing
}
