package botgo

type MessageEmbed struct {
	Title     string                `json:"title"`
	Prompt    string                `json:"prompt"`
	Thumbnail MessageEmbedThumbnail `json:"thumbnail"`
	Fields    []MessageEmbedFiled   `json:"fileds"`
}
