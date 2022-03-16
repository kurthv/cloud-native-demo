# Local Test Run

If you want to show off the power of cloud-native sofware from the privacy of your personal laptop, you can do this too. 

## Requirements

* [KinD: Kubernetes in Docker](https://kind.sigs.k8s.io/docs/user/ingress)
* kubectl
* Helm

## Configure Helm Cluster

```
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF
```

## Install an Ingress Controller

We will choose the NginX Ingress Controller.

Install the controller and then wait for the rollout to complete

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

```

