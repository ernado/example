package entdb

import (
	"github.com/ernado/example/internal/ent"
)

type DB struct {
	ent *ent.Client
}

func New(ent *ent.Client) *DB {
	return &DB{
		ent: ent,
	}
}
