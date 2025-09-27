package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type TelegramChannel struct {
	ent.Schema
}

// Fields of the TelegramChannel.
func (TelegramChannel) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("access_hash"),
		field.String("title"),
		field.Bool("active"),
	}
}
