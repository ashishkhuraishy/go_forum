package token

import "time"

type Maker interface {
	// CreateToken will take in `userID` and `userName`
	// then convert it into a payload, then generates a
	// `token` which expires on the given `duration`
	CreateToken(userId int64, username string, duration time.Duration) (string, error)
	// Verify token method will take in a token and checks
	// it for validity. If the token is valid then the data
	// is converted to a `Payload` and returned else the
	// error will be returned
	VerifyToken(token string) (*Payload, error)
}
