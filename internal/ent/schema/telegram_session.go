package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type TelegramSession struct {
	ent.Schema
}

func (TelegramSession) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.Bytes("data"),
	}
}

func (TelegramSession) Edges() []ent.Edge {
	return []ent.Edge{}
}
