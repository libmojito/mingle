package convo

type Message struct {
	Role     Role              `json:"role" binding:"required"`
	Content  string            `json:"content" binding:"required"`
	FieldIDs []string          `json:"field_ids"`
	Metadata map[string]string `json:"metadata"`
}
