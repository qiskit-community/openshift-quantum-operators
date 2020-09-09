package controller

import (
	"github.com/example-inc/ibm-quantum-operator/pkg/controller/ibmq"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, ibmq.Add)
}
