package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Content       string     `json:"content"`
	ParentID      *uuid.UUID `json:"parent_id,omitempty"` // parent_id can be null
	LikesCount    int        `json:"likes_count"`
	RetweetsCount int        `json:"retweets_count"`
	CommentsCount int        `json:"comments_count"`
	CreatedAt     time.Time  `json:"created_at"`
}
