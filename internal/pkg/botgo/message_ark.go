package botgo

type MessageArk struct {
	TemplateId int            `json:"template_id"`
	Kv         []MessageArkKv `json:"kv"`
}
