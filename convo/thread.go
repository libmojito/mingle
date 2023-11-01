package convo

type Thread struct {
	ID        string            `json:"id"`
	Object    string            `json:"object"`
	CreatedAt int               `json:"created_at"`
	metadata  map[string]string `json:"metadata"`
}
