package sql

import "embed"

//go:embed *.sql
var SchemaFiles embed.FS
