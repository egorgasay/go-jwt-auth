package service

import "fmt"

var (
	ErrNoToken       = fmt.Errorf("no token")
	ErrInvalidToken  = fmt.Errorf("invalid token")
	ErrSign          = fmt.Errorf("can't sign token")
	ErrGenerateUUID  = fmt.Errorf("can't generate uuid")
	ErrExpired       = fmt.Errorf("token expired")
	ErrInvalidGUID   = fmt.Errorf("invalid guid")
	ErrCantHashToken = fmt.Errorf("can't hash token")
)
