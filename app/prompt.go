package app

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/jamesstocktonj1/mcp-provider/bindings/jamesstocktonj1/mcp/mcp_prompt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.wasmcloud.dev/provider"
	wrpc "wrpc.io/go"
	wrpcnats "wrpc.io/go/nats"
)

const TypePrompt = "mcp-prompt"

type PromptHandler interface {
	handlePutTargetLink(link provider.InterfaceLinkDefinition) error
	handleDelTargetLink(link provider.InterfaceLinkDefinition) error
}

type PromptLinker interface {
	AddPrompt(p *mcp.Prompt, h mcp.PromptHandler)
}

type OutgoingClient interface {
	OutgoingRpcClient(target string) *wrpcnats.Client
}

type WrpcPromptHandler func(context.Context, wrpc.Invoker, *mcp_prompt.PromptRequest) (r0__ *wrpc.Result[mcp_prompt.PromptResponse, mcp_prompt.Error], err__ error)

type promptHandler struct {
	links map[string]string

	logger      *slog.Logger
	wrpcHandler WrpcPromptHandler
	client      OutgoingClient

	server PromptLinker
}

func NewPromptHandler(client OutgoingClient, handler PromptLinker, logger *slog.Logger) (*promptHandler, error) {
	return &promptHandler{
		links:       make(map[string]string),
		logger:      logger,
		wrpcHandler: mcp_prompt.Handle,
		client:      client,
		server:      handler,
	}, nil
}

func (h *promptHandler) handlePutTargetLink(link provider.InterfaceLinkDefinition) error {
	args := &mcp.Prompt{}

	name, ok := link.TargetConfig["name"]
	if !ok {
		return errors.New("target config \"name\" is required")
	}
	args.Name = name

	title, ok := link.TargetConfig["title"]
	if ok {
		args.Title = title
	}

	description, ok := link.TargetConfig["description"]
	if ok {
		args.Description = description
	}

	jsonArgs, ok := link.TargetConfig["arguments"]
	if ok {
		promptArgs := []*mcp.PromptArgument{}
		err := json.Unmarshal([]byte(jsonArgs), &promptArgs)
		if err != nil {
			return err
		}
		args.Arguments = promptArgs
	}

	handler, err := h.newPromptHandlerFunc(link.SourceID)
	if err != nil {
		return err
	}

	h.logger.Info("adding prompt", "target", link.SourceID, "args", args, "handler", handler)

	h.server.AddPrompt(args, handler)
	return nil
}

func (h *promptHandler) handleDelTargetLink(link provider.InterfaceLinkDefinition) error {
	h.logger.Info("handling del link", "link", link)
	return nil
}

func (h *promptHandler) newPromptHandlerFunc(target string) (mcp.PromptHandler, error) {
	functionInvoker := h.client.OutgoingRpcClient(target)
	return func(ctx context.Context, r *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		h.logger.Info("new prompt request", "request", r.Session.ID(), "target", target)
		req := mapPromptRequest(r)
		resp, err := h.wrpcHandler(ctx, functionInvoker, req)
		text, ok := resp.Ok.Messages[0].Content.GetText()
		h.logger.Info("prompt responst", "response", text, "present", ok, "description", resp.Ok.Description)
		if err != nil {
			return nil, err
		} else if resp.Err != nil {
			return nil, errors.New(*resp.Err)
		}
		return mapPromptResponse(resp.Ok), nil
	}, nil
}
