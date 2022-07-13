package types

type Common struct {
	Acknowledged bool   `json:"_acknowledged"`
	Error        string `json:"_error,omitempty"`
}
