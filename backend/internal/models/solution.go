package models

type InputField struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Type        string `json:"type"` // "string" | "number" | "boolean"
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
}

type Solution struct {
	ID                string   `json:"id"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	ExpectedBehaviors []string `json:"expected_behaviors"`
	Code              string   `json:"code"`

	// Optional: Information for the UI to render an interactive test form
	Endpoint      string       `json:"endpoint,omitempty"`
	Method        string       `json:"method,omitempty"`
	InputFields   []InputField `json:"input_fields,omitempty"`
	SamplePayload string       `json:"sample_payload,omitempty"`
}
