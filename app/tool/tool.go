package tool

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/jamesstocktonj1/mcp-provider/bindings/jamesstocktonj1/mcp/mcp_prompt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.wasmcloud.dev/provider"
	wrpc "wrpc.io/go"
	wrpcnats "wrpc.io/go/nats"
)

type ToolHandler interface {
	HandlePutTargetLink(link provider.InterfaceLinkDefinition) error
	HandleDelTargetLink(link provider.InterfaceLinkDefinition) error
}

type ToolLinker interface {
	AddTool(t *mcp.Tool, h mcp.ToolHandlerFor[json.RawMessage, json.RawMessage])
}

type OutgoingClient interface {
	OutgoingRpcClient(target string) *wrpcnats.Client
}

type WrpcPromptHandler func(context.Context, wrpc.Invoker, *mcp_prompt.PromptRequest) (r0__ *wrpc.Result[mcp_prompt.PromptResponse, mcp_prompt.Error], err__ error)

type toolHandler struct {
	links map[string]string

	logger      *slog.Logger
	wrpcHandler WrpcPromptHandler
	client      OutgoingClient

	server ToolLinker
}

func NewToolHandler(client OutgoingClient, server *mcp.Server, logger *slog.Logger) (*toolHandler, error) {
	return nil, nil
}

func (h *toolHandler) HandlePutTargetLink(link provider.InterfaceLinkDefinition) error {
	return nil
}

func (h *toolHandler) HandleDelTargetLink(link provider.InterfaceLinkDefinition) error {
	return nil
}
