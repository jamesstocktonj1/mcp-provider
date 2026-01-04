package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jamesstocktonj1/mcp-provider/app/prompt"
	"github.com/jamesstocktonj1/mcp-provider/app/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.wasmcloud.dev/provider"
)

type server struct {
	provider *provider.WasmcloudProvider

	mcpServer  *mcp.Server
	httpServer *http.Server

	errorChan  chan error
	signalChan chan os.Signal

	// wRPC handlers
	promptHandler prompt.PromptHandler
	toolHandler   tool.ToolHandler
}

type args struct {
	Name string `json:"name" jsonschema:"the person to greet"`
}

func NewServer() (*server, error) {
	s := &server{}

	// Create new Wasmcloud Provider
	prov, err := provider.New(
		provider.TargetLinkPut(s.handlePutTargetLink),
		provider.TargetLinkDel(s.handleDelTargetLink),
		provider.HealthCheck(s.handleHealth),
		provider.Shutdown(s.handleShutdown),
	)
	if err != nil {
		return nil, err
	}
	s.provider = prov

	// Create new MCP Server
	s.mcpServer = mcp.NewServer(&mcp.Implementation{
		Name: "greeter",
	}, nil)
	s.mcpServer.AddReceivingMiddleware(s.loggingMiddleware)
	s.mcpServer.AddSendingMiddleware(s.loggingMiddleware)

	// Create new wRPC handlers
	s.promptHandler, err = prompt.NewPromptHandler(prov, s.mcpServer, prov.Logger)
	if err != nil {
		return nil, err
	}

	s.toolHandler, err = tool.NewToolHandler(prov, s.mcpServer, prov.Logger)
	if err != nil {
		return nil, err
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "greet",
		Description: "say hi",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args args) (*mcp.CallToolResult, any, error) {
		s.provider.Logger.Info("Greeting Someone", "params", req.GetParams())
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Hi %+v", args)},
			},
		}, nil, nil
	})

	s.mcpServer.AddResource(&mcp.Resource{
		URI:      "file:///project/messages.txt",
		MIMEType: "text/plain",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		s.provider.Logger.Info("Read Messages", "params", req.GetParams())
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:  "file:///project/messages.txt",
					Text: "Hello Messages",
				},
			},
		}, nil
	})

	// Create new HTTP Server
	s.httpServer = &http.Server{
		Addr: ":8080",
		Handler: mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			s.provider.Logger.Info("New Request", "host", r.URL.Host, "path", r.URL.Path)
			return s.mcpServer
		}, nil),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return s, nil
}

func (s *server) Run() error {
	// Setup shutdown channels
	s.errorChan = make(chan error, 1)
	s.signalChan = make(chan os.Signal, 1)

	signal.Notify(s.signalChan, syscall.SIGINT)

	// Start Provider
	go func() {
		err := s.provider.Start()
		s.errorChan <- err
	}()

	// Start MCP Http Server
	go func() {
		err := s.httpServer.ListenAndServe()
		s.errorChan <- err
	}()

	// Shutdown after return
	defer func() {
		if err := s.Stop(); err != nil {
			log.Printf("unable to stop server - %s", err.Error())
		}
	}()

	// Run until shutdown is detected
	select {
	case err := <-s.errorChan:
		return err
	case <-s.signalChan:
		return errors.New("server exited with SIGINT")
	}
}

func (s *server) Stop() error {
	if err := s.httpServer.Close(); err != nil {
		return err
	}

	if err := s.provider.Shutdown(); err != nil {
		return err
	}
	return nil
}
