package constants

import "fmt"

var (
	ErrMissingRefreshToken = fmt.Errorf("refresh token was not provided")
	ErrMissingAccessToken  = fmt.Errorf("access token was not provided")
	ErrInvalidToken        = fmt.Errorf("invalid token")
	ErrSignToken           = fmt.Errorf("can't sign token")
	ErrGenerateToken       = fmt.Errorf("can't generate token")
	ErrTokenExpired        = fmt.Errorf("token expired")
	ErrInvalidGUID         = fmt.Errorf("invalid guid")
	ErrCantHashToken       = fmt.Errorf("can't hash token")
)
