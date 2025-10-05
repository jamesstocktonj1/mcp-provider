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
		fmt.Printf("Greeting Someone - %+v\n", req)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Hi %+v", args)},
			},
		}, nil, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name: "greetings",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "name",
				Description: "name input",
				Required:    true,
			},
		},
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		fmt.Printf("Prompt Someone - %+v\n", req.GetParams())
		return &mcp.GetPromptResult{
			Description: "this is some greeting",
			Messages: []*mcp.PromptMessage{
				{
					Content: &mcp.TextContent{Text: "Hello Some Prompt"},
					Role:    mcp.Role("user"),
				},
			},
		}, nil
	})

	s.mcpServer.AddResource(&mcp.Resource{
		URI:      "file:///project/messages.txt",
		MIMEType: "text/plain",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		fmt.Printf("Read Messages - %+v\n", req.GetParams())
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:  "file:///project/messages.txt",
					Text: "Hello Messages",
				},
			},
		}, nil
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
