package user_settings

import "github.com/Traliaa/KineticVPN-Bot/internal/pg/user_settings/sql"

// UserSettings implement db store
type UserSettings struct {
	sql *sql.Queries
}

// New instance
func New() *UserSettings {
	return &UserSettings{
		sql: sql.New(),
	}
}
