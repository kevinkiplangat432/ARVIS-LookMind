package proxy

import "encoding/json"

// usageResponse captures just the fields ARVIS needs out of a
// provider's response body. "usage.prompt_tokens" / "usage.
// completion_tokens" is the de facto standard OpenAI set, and most
// OpenAI-compatible providers mirror this shape.
type usageResponse struct {
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}

// extractTokenUsage never errors — an unparseable or missing usage
// field just means zero counts logged, not a failed request. Token
// accounting is a nice-to-have on top of the audit trail, not a
// gate on whether the trail gets written at all.
func extractTokenUsage(body []byte) (promptTokens, completionTokens int) {
	var parsed usageResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, 0
	}
	return parsed.Usage.PromptTokens, parsed.Usage.CompletionTokens
}