package app

import (
	"context"

	"github.com/jamesstocktonj1/mcp-provider/bindings/jamesstocktonj1/mcp/mcp_prompt"
	"go.wasmcloud.dev/provider"
	wrpc "wrpc.io/go"
)

const TypePrompt = "mcp-prompt"

type PromptHandler interface {
	handlePutTargetLink(link provider.InterfaceLinkDefinition) error
	handleDelTargetLink(link provider.InterfaceLinkDefinition) error
}

type WrpcPromptHandler func(context.Context, wrpc.Invoker, *mcp_prompt.PromptRequest) (r0__ *wrpc.Result[mcp_prompt.PromptResponse, mcp_prompt.Error], err__ error)

type promptHandler struct {
	links map[string]string

	wrpcHandler WrpcPromptHandler
}

func NewPromptHandler() (*promptHandler, error) {
	return &promptHandler{
		links:       make(map[string]string),
		wrpcHandler: mcp_prompt.Handle,
	}, nil
}

func (h *promptHandler) handlePutTargetLink(link provider.InterfaceLinkDefinition) error {
	return nil
}

func (h *promptHandler) handleDelTargetLink(link provider.InterfaceLinkDefinition) error {
	return nil
}
