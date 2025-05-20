package common

import "errors"

var (
	ErrRedirected = errors.New("The client has been redirected")
)
