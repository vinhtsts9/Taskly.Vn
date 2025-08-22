package embedSchema

import "embed"

//go:embed schema/*.sql
var Files embed.FS