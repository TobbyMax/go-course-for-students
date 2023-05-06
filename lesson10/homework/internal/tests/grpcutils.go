package tests

import "errors"

var (
	ErrUserNotFound    = errors.New("rpc error: code = NotFound desc = user with such id does not exist")
	ErrAdNotFound      = errors.New("rpc error: code = NotFound desc = ad with such id does not exist")
	ErrGRPCForbidden   = errors.New("rpc error: code = PermissionDenied desc = forbidden")
	ErrInvalidEmail    = errors.New("rpc error: code = InvalidArgument desc = mail: missing '@' or angle-addr")
	ErrMissingArgument = errors.New("rpc error: code = InvalidArgument desc = required argument is missing")
	ErrMockInternal    = errors.New("rpc error: code = Internal desc = mock error")
	ErrValidationMock  = errors.New("rpc error: code = InvalidArgument desc = ")
	ErrDateMock        = errors.New("rpc error: code = InvalidArgument desc = parsing time \"20/02/2022\" as \"2006-01-02\": cannot parse \"2/2022\" as \"2006\"")
)
