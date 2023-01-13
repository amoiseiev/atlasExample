package dbschema

import "embed"

//go:embed *.sql
var SQLFiles embed.FS
