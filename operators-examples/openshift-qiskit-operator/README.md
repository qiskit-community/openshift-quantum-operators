# Overview

## Artifacts Detail

    .
        ├── image                   # Files required for building Jupyter Notebook image
        ├── operator                # Operator artifacts and documentation 
        └── README.md

## Usage

### Installation from Image on OpenShift 4

This installation method will use the latest version of the operator image that has been built and published to Quay

#### Deploy the custom resource definition (CRD)

```
oc apply -f deploy/crds/singhp11.io_qiskitplaygrounds_crd.yaml
```

#### Deploy the RBAC configuration:

```
oc apply -f deploy/role.yaml
oc apply -f deploy/service_account.yaml
oc apply -f deploy/role_binding.yaml
```

#### Setting up authorization with IBMQ Account token

- Edit the configuration file:

```
deploy/secret.cfg
[AUTH TOKENS]
token = your_IBMQ_account_token
```

- Convert the configuration to base64:

```
cat secret.cfg | base64
```

- Place the output in deploy/secret.yaml as:

```
apiVersion: v1
kind: Secret
metadata:
	name: qiskit-secret
type: Opaque
data:
	qiskit-secret.cfg: <base64 encoded secret.cfg>
```

- Deploy the secret

```
oc apply -f deploy/secret.yaml
```

#### Deploy the operator itself

```
oc apply -f deploy/operator.yaml
```

#### Wait for the operator pod deployment to complete

```
watch oc get pods
```

#### Deploy an instance of the custom resource

```
oc apply -f deploy/crds/singhp11.io_v1_qiskitplayground_cr.yaml
```

#### The notebook is found on the exposed route

```
oc get routes
```
