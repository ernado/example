package db

import (
	"embed"
)

//go:embed _migrations/*.sql
var Migrations embed.FS
