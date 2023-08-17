package botgo

type Member struct {
	User     User      `json:"user"`
	Nick     string    `json:"nick"`
	Roles    []string  `json:"roles"`
	JoinedAt Timestamp `json:"joined_at"`
}
