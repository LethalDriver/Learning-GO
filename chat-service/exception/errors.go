package exception

import "errors"

var ErrRoomExists = errors.New("room already exists")
var ErrEntityNotFound = errors.New("not found")
