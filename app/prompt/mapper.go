package prompt

import (
	"github.com/jamesstocktonj1/mcp-provider/bindings/jamesstocktonj1/mcp/mcp_prompt"
	mcp_types "github.com/jamesstocktonj1/mcp-provider/bindings/jamesstocktonj1/mcp/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	wrpc "wrpc.io/go"
)

func mapPromptRequest(req *mcp.GetPromptRequest) *mcp_prompt.PromptRequest {
	return &mcp_prompt.PromptRequest{
		Session: mapSession(req.Session),
		Params:  mapPromptParams(req.Params),
	}
}

func mapPromptResponse(res *mcp_prompt.PromptResponse) *mcp.GetPromptResult {
	return &mcp.GetPromptResult{
		Description: res.Description,
		Messages:    mapMessages(res.Messages),
		Meta:        mcp.Meta{},
	}
}

func mapSession(s *mcp.ServerSession) *mcp_types.Session {
	return &mcp_types.Session{
		Id: s.ID(),
	}
}

func mapPromptParams(v *mcp.GetPromptParams) *mcp_prompt.Params {
	return &mcp_prompt.Params{
		Name:      v.Name,
		Arguments: mapMap(v.Arguments),
	}
}

func mapMap[U comparable, V any](u map[U]V) []*wrpc.Tuple2[U, V] {
	v := make([]*wrpc.Tuple2[U, V], len(u))
	i := 0
	for x, y := range u {
		v[i] = &wrpc.Tuple2[U, V]{
			V0: x,
			V1: y,
		}
		i++
	}
	return v
}

func mapMessages(msgs []*mcp_prompt.PromptMessage) []*mcp.PromptMessage {
	m := make([]*mcp.PromptMessage, len(msgs))
	for i, msg := range msgs {
		m[i] = mapMessage(msg)
	}
	return m
}

func mapMessage(msg *mcp_prompt.PromptMessage) *mcp.PromptMessage {
	return &mcp.PromptMessage{
		Content: mapContent(msg.Content),
		Role:    mapRole(msg.Role),
	}
}

func mapContent(c *mcp_types.Content) mcp.Content {
	switch c.Discriminant() {
	case mcp_types.ContentText:
		txt, _ := c.GetText()
		return &mcp.TextContent{
			Text:        txt.Text,
			Meta:        mcp.Meta{},
			Annotations: &mcp.Annotations{},
		}
	default:
		return nil
	}
}

func mapRole(r mcp_prompt.Role) mcp.Role {
	switch r.String() {
	case "user":
		return mcp.Role("user")
	default:
		return mcp.Role("")
	}
}
