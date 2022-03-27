# cloud-native-demo
A demo for the power of cloud native application design


## Installation

Requirements:
* a Kubernetes Cluster to which you have writing API access.
* the [step-cli](https://github.com/smallstep/cli) locally installed.

TODO: make this into a shell script and execute it via a docker image.

```
step certificate create root.linkerd.cluster.local ca.crt ca.key \
--profile root-ca --no-password --insecure

step certificate create identity.linkerd.cluster.local issuer.crt issuer.key \
--profile intermediate-ca --not-after 8760h --no-password --insecure \
--ca ca.crt --ca-key ca.key

helm install cloud-native-demo \
  --set-file linkerd2.identityTrustAnchorsPEM=ca.crt \
  --set-file linkerd2.identity.issuer.tls.crtPEM=issuer.crt \
  --set-file linkerd2.identity.issuer.tls.keyPEM=issuer.key \
  cloud-native-demo
```

## Quick Check

```
kubectl proxy

```
Web-Browser to: 

* Bug Simulator: http://127.0.0.1:8001/api/v1/namespaces/cloud-native-demo/services/cloud-native-demo:80/proxy/
* Emojivoto Web Frontend: http://127.0.0.1:8001/api/v1/namespaces/cloud-native-demo/services/cloud-native-demo-emojivoto-web-svc:80/proxy/