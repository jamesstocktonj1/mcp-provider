package app

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.wasmcloud.dev/provider"
)

func (s *server) handlePutTargetLink(link provider.InterfaceLinkDefinition) error {
	s.provider.Logger.Info("Handling Put Link", "link", link)
	return nil
}

func (s *server) handleDelTargetLink(link provider.InterfaceLinkDefinition) error {
	s.provider.Logger.Info("Handling Delete Link", "link", link)
	return nil
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
