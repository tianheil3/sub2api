package xai

// Model describes an xAI model in OpenAI-compatible /models shape.
type Model struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	Created     int64  `json:"created,omitempty"`
	OwnedBy     string `json:"owned_by"`
	DisplayName string `json:"display_name,omitempty"`
}

const (
	DefaultCLIClientVersion    = "0.2.22"
	DefaultCLIClientIdentifier = "grok-pager"
	DefaultCLITokenAuth        = "xai-grok-cli"
)

var defaultModels = []Model{
	{ID: "grok-4.3", Object: "model", OwnedBy: "xai", DisplayName: "Grok 4.3"},
	{ID: "grok-composer-2.5-fast", Object: "model", OwnedBy: "xai", DisplayName: "Composer 2.5 Fast"},
	{ID: "grok-build-0.1", Object: "model", OwnedBy: "xai", DisplayName: "Grok Build 0.1"},
	{ID: "grok-4.20-0309-reasoning", Object: "model", OwnedBy: "xai", DisplayName: "Grok 4.20 Reasoning"},
	{ID: "grok-4.20-0309-non-reasoning", Object: "model", OwnedBy: "xai", DisplayName: "Grok 4.20 Non Reasoning"},
	{ID: "grok-4.20-multi-agent-0309", Object: "model", OwnedBy: "xai", DisplayName: "Grok 4.20 Multi Agent"},
}

func DefaultModels() []Model {
	out := make([]Model, len(defaultModels))
	copy(out, defaultModels)
	return out
}

func DefaultModelIDs() []string {
	models := DefaultModels()
	ids := make([]string, 0, len(models))
	for _, model := range models {
		ids = append(ids, model.ID)
	}
	return ids
}

func DefaultModelMapping() map[string]string {
	mapping := make(map[string]string, len(defaultModels)+3)
	for _, model := range defaultModels {
		mapping[model.ID] = model.ID
	}
	mapping["grok"] = "grok-4.3"
	mapping["grok-latest"] = "grok-4.3"
	mapping["grok-build"] = "grok-build-0.1"
	mapping["composer-2.5-fast"] = "grok-composer-2.5-fast"
	mapping["composer 2.5 fast"] = "grok-composer-2.5-fast"
	mapping["Composer 2.5 Fast"] = "grok-composer-2.5-fast"
	mapping["grok-4.20-reasoning"] = "grok-4.20-0309-reasoning"
	mapping["grok-4.20-non-reasoning"] = "grok-4.20-0309-non-reasoning"
	return mapping
}

func IsCLIModel(model string) bool {
	switch model {
	case "grok-build", "grok-build-0.1", "grok-composer-2.5-fast":
		return true
	default:
		return false
	}
}

func BaseURLForModel(accountBaseURL string, model string) string {
	if IsCLIModel(model) {
		return DefaultCLIBaseURL
	}
	return accountBaseURL
}

func CLIHeaders(model string) map[string]string {
	if !IsCLIModel(model) {
		return nil
	}
	userAgent := "grok-pager/" + DefaultCLIClientVersion + " grok-shell/" + DefaultCLIClientVersion + " (macos; aarch64)"
	return map[string]string{
		"User-Agent":               userAgent,
		"x-grok-client-identifier": DefaultCLIClientIdentifier,
		"x-grok-client-version":    DefaultCLIClientVersion,
		"x-xai-token-auth":         DefaultCLITokenAuth,
		"x-grok-model-override":    model,
	}
}

func ApplyCLIHeaders(headers interface{ Set(string, string) }, model string) {
	for key, value := range CLIHeaders(model) {
		headers.Set(key, value)
	}
}
