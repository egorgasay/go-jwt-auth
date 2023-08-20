package usecase

import "fmt"

var (
	ErrNoToken      = fmt.Errorf("no token")
	ErrInvalidToken = fmt.Errorf("invalid token")
	ErrSign         = fmt.Errorf("can't sign token")
	ErrGenerateUUID = fmt.Errorf("can't generate uuid")
)
