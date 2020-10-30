package controller

import (
	"github.com/youaij/operator-demo/pkg/controller/learn"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, learn.Add)
}
