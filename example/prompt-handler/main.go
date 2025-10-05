package main

import (
	mcpprompt "github.com/jamesstocktonj1/mcp-handler/example/prompt-handler/gen/jamesstocktonj1/mcp/mcp-prompt"
	"github.com/jamesstocktonj1/mcp-handler/example/prompt-handler/gen/jamesstocktonj1/mcp/types"
	"go.bytecodealliance.org/cm"
)

func init() {
	mcpprompt.Exports.Handle = PromptHandler
}

func PromptHandler(request mcpprompt.PromptRequest) (result cm.Result[mcpprompt.PromptResponseShape, mcpprompt.PromptResponse, mcpprompt.Error]) {
	return cm.OK[cm.Result[mcpprompt.PromptResponseShape, mcpprompt.PromptResponse, mcpprompt.Error]](mcpprompt.PromptResponse{
		Description: "this is some greeting",
		Messages: cm.ToList([]mcpprompt.PromptMessage{
			{
				Role: types.RoleUser,
				Content: types.ContentText(types.TextContent{
					Text: "Hello Some Prompt",
				}),
			},
		}),
	})
}

//go:generate go tool wit-bindgen-go generate --world component --out gen ./wit
func main() {}
