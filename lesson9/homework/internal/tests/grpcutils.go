package tests

import "errors"

var (
	ErrUserNotFound  = errors.New("rpc error: code = NotFound desc = user with such id does not exist")
	ErrAdNotFound    = errors.New("rpc error: code = NotFound desc = ad with such id does not exist")
	ErrGRPCForbidden = errors.New("rpc error: code = PermissionDenied desc = forbidden")
	ErrInvalidEmail  = errors.New("rpc error: code = InvalidArgument desc = mail: missing '@' or angle-addr")
)
