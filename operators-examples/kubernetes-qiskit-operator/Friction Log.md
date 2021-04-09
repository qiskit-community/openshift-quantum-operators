The Qiskit Playground Operator is a foundation of self-service Qiskit deployments on Kubernetes.

The original work was done by the Red Hat Office of the CTO and now it is continued at IBM Quantum TECS team in partnership with Red Hat.

This version is based on a current Operator SDK.


#Prerequisutes

Have Go installed
Kubernetes environment configured for _kubectl_

#Deployment

check out from _git_

```
make install run
```

Note: need to put in the background if run with make from a terminal

The operator defines a custom resource, CRD, QiskitPlayground.

To create an instance, CR, define it in a yaml file, test.yaml:

```
kind: QiskitPlayground
metadata:
  name: test
spec:
  image: "jupyter/scipy-notebook:latest"
```

and apply as usual with

```
kubectl apply -f test.yaml
```

#Exposing the Service

On minikube, do

```
minikube service --url test-service
```

to obtain the URL for the service to open in a browser.

Now you should have a Jupyter notebook with qiskit in your browser.

Note: we need to actually populate it with qiskit and some notebooks.
