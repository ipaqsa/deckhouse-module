package main

import (
	"context"

	"github.com/x/module/internal/config"
	"github.com/x/module/internal/server"
)

const address = ":8080"

func main() {
	if err := server.New(config.MustLoad()).Serve(context.Background(), address); err != nil {
		panic(err)
	}
}
