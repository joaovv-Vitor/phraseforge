// Package migrations embeds the SQL migrations applied by the application.
package migrations

import "embed"

// Files contains all SQL migration files.
//
//go:embed *.sql
var Files embed.FS
