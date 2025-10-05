package main

import (
	"log"

	"github.com/jamesstocktonj1/mcp-provider/app"
)

//go:generate wit-bindgen-wrpc go --out-dir bindings --world imports --package github.com/jamesstocktonj1/mcp-provider/bindings wit
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	s, err := app.NewServer()
	if err != nil {
		return err
	}

	return s.Run()
}
