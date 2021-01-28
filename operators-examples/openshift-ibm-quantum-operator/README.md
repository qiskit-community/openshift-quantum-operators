# Overview

## Artifacts Detail

    .
        ├── image                   # Files required for building cr image
        ├── operator                # Operator artifacts and documentation 
        └── README.md

## Introduction

OpenShift IBM Quantum Operator creates a flexible serving system to build quantum circuits using Qiskit.

Here are some key features:

### Flexible REST endpoints for running quantum workloads in a Kubernetes cluster

You can submit quantum workloads written in Qiskit which are executed on IBM Quantum systems and simulators. Workloads are executed as pods, orchestrated and managed by Kube APIs  

### Environment Configuration & Scheduling

The operator authenticates the caller and mounts the equivalent secrets on the quantum workloads pods. The default scheduler picks the least busy backend that has enough qubits to run the workload.

### High Availability via Horizontal Pod Autoscaling

Run multiple instances of the server system for high availability. Either specify
a static number of replicas or easily configure horizontal auto scaling to create (and delete)
instances based on resource consumption.

### IBM Quantum Experience and Qiskit

IBM Quantum Experience is quantum on the cloud. With IBM Quantum Experience you can:

- Run - Run quantum circuits on IBM Quantum systems and simulators. Users in the IBM Quantum Network can also access premium systems.
- Develop and deploy - Explore quantum applications in areas such as chemistry, optimization, finance, and AI.
- Quantum innovation - Stay informed and contribute to the future of quantum. 

Qiskit is an open-source framework for working with quantum computers at the level of circuits, pulses, and algorithms. A central goal of Qiskit is
to build a software stack that makes it easy for anyone to use quantum computers.

#### Note

The user needs to ensure that their IBM Quantum API token / secret exists in their Kuberenetes namespace/project where the operator will be installed.
Any api could be changed any time without any notice. That said, your feedback is very important and appreciated to make this project more stable and useful.

### Contributing

You can contribute by

* Raising any [issues](https://github.com/qiskit-community/openshift-quantum-operators/issues) you find using openshift-ibm-quantum-operator community operator
* Fixing issues by opening [Pull Requests](https://github.com/qiskit-community/openshift-quantum-operators/pulls)
* Talking about openshift-ibm-quantum-operator
  
### License

openshift-ibm-quantum-operator is licensed under the [Apache 2.0 license](https://github.com/qiskit-community/openshift-quantum-operators/blob/master/LICENSE)

## Usage

### Installation from Image on OpenShift 4

This installation method will use the latest version of the operator image that has been built and published to Quay

#### Deploy the custom resource definition (CRD)

``` bash
oc apply -f deploy/crds/singhp11.io_ibmqes_crd.yaml

```

#### Deploy the RBAC configuration

``` bash
oc apply -f deploy/role.yaml
oc apply -f deploy/service_account.yaml
oc apply -f deploy/role_binding.yaml
```

#### Setting up authorization with an IBM Quantum account token

- Edit the configuration file:

``` bash
deploy/secret.cfg
[AUTH TOKENS]
token = your_IBM_Quantum_account_token
```

- Convert the configuration to base64:

``` bash
cat secret.cfg | base64
```

- Place the output in deploy/secret.yaml as:

``` bash
apiVersion: v1
kind: Secret
metadata:
	name: qiskit-secret
type: Opaque
data:
	qiskit-secret.cfg: <base64 encoded secret.cfg>
```

- Deploy the secret

``` bash
oc apply -f deploy/secret.yaml

```

#### Deploy the operator itself

``` bash
oc apply -f deploy/operator.yaml
```

#### Wait for the operator pod deployment to complete

``` bash
watch oc get pods
```

#### Deploy an instance of the custom resource

``` bash
oc apply -f deploy/crds/singhp11.io_v1_ibmqe_cr.yaml
```

#### Access the serving system via the routes

``` bash
oc get routes
```
