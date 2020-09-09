# Usage

  

## Installation from Image on OpenShift 4

  

This installation method will use the latest version of the operator image that has been built and published to docker hub

  

1. Deploy the custom resource definition (CRD):

```
oc apply -f deploy/crds/dobtech.io_qiskitplaygrounds_crd.yaml

```

 2. Deploy the RBAC configuration:
```
oc apply -f deploy/role.yaml
oc apply -f deploy/service_account.yaml
oc apply -f deploy/role_binding.yaml
```
 3. Setting up authorization with IBMQ Account token
 
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
4. Deploy the operator itself:

```
oc apply -f deploy/operator.yaml
```
5. Wait for the operator pod deployment to complete:
```
watch oc get pods
```

6. Deploy an instance of the custom resource:
```
oc apply -f deploy/crds/dobtech.io_v1_qiskitplayground_cr.yaml
```
7. The notebook is found on the exposed route