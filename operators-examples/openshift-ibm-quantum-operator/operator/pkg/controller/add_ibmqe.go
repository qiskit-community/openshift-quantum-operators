package controller

import (
	"github.com/example-inc/openshift-ibm-quantum-operator/pkg/controller/ibmqe"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, ibmqe.Add)
}
