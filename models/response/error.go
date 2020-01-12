package response

import "errors"

var (
	ErrNotFound        = errors.New(`Not Found`)
	ErrServer          = errors.New(`Internal Server Error`)
	ErrUnAuthorized    = errors.New(`UnAuthorized`)
	ErrForbidden       = errors.New(`Forbidden`)
	ErrAlreadyExist    = errors.New(`Already Exist`)
	ErrBadRequest      = errors.New(`Bad Request`)
	ErrNoDeviceForRoom = errors.New(`No Device Set Up for This Room`)
	ErrLogin           = errors.New(`Invalid Email or Password`)
	ErrInterface       = errors.New(`Service Unavailable`)
)
