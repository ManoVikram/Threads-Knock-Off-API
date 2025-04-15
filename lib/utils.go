package lib

import "database/sql"

// Helper function to return an empty string if sql.NullString is invalid
func StringOrEmpty(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
