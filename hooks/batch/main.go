package main

import (
	"github.com/deckhouse/module-sdk/pkg/app"

	_ "hooks/hooks"
)

func main() {
	app.Run()
}
