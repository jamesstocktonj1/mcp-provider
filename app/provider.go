package app

import (
	"context"
	"fmt"

	"github.com/jamesstocktonj1/mcp-provider/app/prompt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.wasmcloud.dev/provider"
)

func (s *server) handlePutTargetLink(link provider.InterfaceLinkDefinition) error {
	if len(link.Interfaces) > 1 {
		s.provider.Logger.Warn("multiple link interfaces defined, using first", "link", link)
	}

	switch link.Interfaces[0] {
	case prompt.TypePrompt:
		return s.promptHandler.HandlePutTargetLink(link)
	default:
		s.provider.Logger.Error("unknown interface type - "+link.Interfaces[0], "link", link)
		return fmt.Errorf("unknown interface type - %s", link.Interfaces[0])
	}
}

func (s *server) handleDelTargetLink(link provider.InterfaceLinkDefinition) error {
	if len(link.Interfaces) > 1 {
		s.provider.Logger.Warn("multiple link interfaces defined, using first", "link", link)
	}

	switch link.Interfaces[0] {
	case prompt.TypePrompt:
		return s.promptHandler.HandleDelTargetLink(link)
	default:
		s.provider.Logger.Error("unknown interface type - "+link.Interfaces[0], "link", link)
		return fmt.Errorf("unknown interface type - %s", link.Interfaces[0])
	}
}

func (s *server) handleHealth() string {
	return "healthy"
}

func (s *server) handleShutdown() error {
	return nil
}

func (s *server) loggingMiddleware(next mcp.MethodHandler) mcp.MethodHandler {
	return func(ctx context.Context, method string, req mcp.Request) (result mcp.Result, err error) {
		s.provider.Logger.Info("MCP Middleware", "method", method, "session_id", req.GetSession().ID())
		return next(ctx, method, req)
	}
}
