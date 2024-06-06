package error_service

import "errors"

var ErrNoPhoneNumber = errors.New("phone number doesn't exist")
var ErrNoUser = errors.New("user does not exist")
var ErrDuplicatePhoneNumber = errors.New("phone number already in use")
var ErrBidAcceptance = errors.New("bid could not be accepted")
