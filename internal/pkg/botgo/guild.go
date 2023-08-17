package botgo

type Guild struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	OwnerId     string `json:"owner_id"`
	Owner       bool   `json:"owner"`
	MemberCount int    `json:"member_count"`
	MaxMembers  int    `json:"max_members"`
	Description string `json:"description"`
	JoinedAt    string `json:"joined_at"`
}
