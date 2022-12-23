package golang

import "errors"

var ErrMultipleReceivers = errors.New("expected method to have one receiver")
