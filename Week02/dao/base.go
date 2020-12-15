package dao

import (
	"github.com/pkg/errors"
)

//type ErrNoRowsFounder interface {
//	error
//	IsNoRows() bool
//}


type errNoRows struct {
	Msg string
}

func (e errNoRows) Error() string {
	return e.Msg
}

func (e errNoRows) Is(target error) bool {
	return errors.Is(target, ErrNoRows)
}

var ErrNoRows = errNoRows{"dao : The rows are not found"}
