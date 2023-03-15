package discord

type SendErrorFrom string

const (
	DB      SendErrorFrom = "DB"
	Discord SendErrorFrom = "Discord"
)
