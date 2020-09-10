package controller

import (
	"github.com/husky-parul/openshift-qiskit-operator/pkg/controller/qiskitplayground"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, qiskitplayground.Add)
}
