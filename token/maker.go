package token

import "time"

// manage tokens
type Maker interface {
	// create token for username with duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// check if token valid
	VerifyToken(token string) (*Payload, error)
}
 