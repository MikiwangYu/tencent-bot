package botgo

type Message struct {
	Id               string              `json:"id"`
	ChannelId        string              `json:"channel_id"`
	GuildId          string              `json:"guild_id"`
	Content          string              `json:"content"`
	Timestamp        Timestamp           `json:"timestamp"`
	EditedTimestamp  Timestamp           `json:"edited_timestamp"`
	MentionEveryone  bool                `json:"memtion_everyone"`
	Author           User                `json:"author"`
	Attachments      []MessageAttachment `json:"attachments"`
	Embeds           []MessageEmbed      `json:"embeds"`
	Mentions         []User              `json:"mentions"`
	Member           Member              `json:"member"`
	Ark              MessageArk          `json:"ark"`
	Seq              int                 `json:"seq"`
	SeqInChanel      string              `json:"seq_in_channel"`
	MessageReference MessageReference    `json:"message_reference"`
	SrcGuildId       string              `json:"src_guild_id"`
}
