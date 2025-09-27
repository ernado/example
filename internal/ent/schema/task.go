package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Task struct {
	ent.Schema
}

func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
	}
}
