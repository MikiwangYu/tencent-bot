package botgo

type Payload struct {
	Op int    `json:"op"`
	D  any    `json:"d"`
	S  int    `json:"s"`
	T  string `json:"t"`
}
