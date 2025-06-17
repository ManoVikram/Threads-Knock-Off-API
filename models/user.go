package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	EmailVerified  sql.NullTime   `json:"emailVerified"`
	Image          sql.NullString `json:"image"`
	Username       sql.NullString `json:"username"`
	FollowerCount  int            `json:"followerCount"`
	FollowingCount int            `json:"followingCount"`
	Bio            sql.NullString `json:"bio"`
}
