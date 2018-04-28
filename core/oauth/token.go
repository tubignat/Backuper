package oauth

import (
	"time"
)

// Token contains a string representation and expiration date of OAuth token
type Token struct {
	Value   string
	Expires time.Time
}
