package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:generate wit-bindgen-wrpc go --out-dir bindings --world imports --package github.com/jamesstocktonj1/ticker-provider/bindings wit
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	server := mcp.NewServer(&mcp.Implementation{
		Name: "greeter",
	}, nil)

	type args struct {
		Name string `json:"name" jsonschema:"the person to greet"`
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "say hi",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args args) (*mcp.CallToolResult, any, error) {
		fmt.Printf("Greeting Someone - %s\n", args.Name)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Hi " + args.Name},
			},
		}, nil, nil
	})

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Starting server on address - %s\n", httpServer.Addr)
	return httpServer.ListenAndServe()
}
