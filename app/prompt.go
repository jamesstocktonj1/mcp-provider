package app

import "go.wasmcloud.dev/provider"

const TypePrompt = "mcp-prompt"

type PromptHandler interface {
	handlePutTargetLink(link provider.InterfaceLinkDefinition) error
	handleDelTargetLink(link provider.InterfaceLinkDefinition) error
}

type promptHandler struct {
	links map[string]string
}

func NewPromptHandler() (*promptHandler, error) {
	return &promptHandler{
		links: make(map[string]string),
	}, nil
}

func (h *promptHandler) handlePutTargetLink(link provider.InterfaceLinkDefinition) error {
	return nil
}

func (h *promptHandler) handleDelTargetLink(link provider.InterfaceLinkDefinition) error {
	return nil
}
