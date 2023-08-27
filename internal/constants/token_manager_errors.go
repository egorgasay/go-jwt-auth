package constants

import "fmt"

var (
	ErrNoToken       = fmt.Errorf("the token was not provided")
	ErrInvalidToken  = fmt.Errorf("invalid token")
	ErrSignToken     = fmt.Errorf("can't sign token")
	ErrGenerateUUID  = fmt.Errorf("can't generate uuid")
	ErrTokenExpired  = fmt.Errorf("token expired")
	ErrInvalidGUID   = fmt.Errorf("invalid guid")
	ErrCantHashToken = fmt.Errorf("can't hash token")
)
