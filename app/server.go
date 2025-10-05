package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type server struct {
	mcpServer  *mcp.Server
	httpServer *http.Server
}

type args struct {
	Name string `json:"name" jsonschema:"the person to greet"`
}

func NewServer() (*server, error) {
	s := &server{}

	s.mcpServer = mcp.NewServer(&mcp.Implementation{
		Name: "greeter",
	}, nil)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
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

	s.httpServer = &http.Server{
		Addr: ":8080",
		Handler: mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			return s.mcpServer
		}, nil),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return s, nil
}

func (s *server) Run() error {
	return s.httpServer.ListenAndServe()
}
