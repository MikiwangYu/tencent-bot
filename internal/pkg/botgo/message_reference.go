package botgo

type MessageReference struct {
	MessageId             string `json:"message_id"`
	IgnoreGetMessageError bool   `json:"ignore_get_message_error"`
}
