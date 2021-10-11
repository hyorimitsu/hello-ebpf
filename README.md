Hello eBPF
---

This is a sample of running eBPF program in kubernetes.

## Description

This program is a sample to run the eBPF program (by BPF Compiler Collection) in minikube.
The configuration is as shown in the figure below.
<!-- TODO: 構成図 -->

## Directory Structure

```
.
├── app-api      # => application source
├── bcc          # => eBPF sources
│    └── tools   # => custom eBPF tools
└── k8s          # => k8s definitions
```

## Usage

### 1. Run the application

#### 1.1 With eBPF programs

a. start minikube

```shell
minikube start --iso-url https://storage.googleapis.com/minikube-performance/minikube.iso --driver=virtualbox
```

b. download and extract necessary kernel headers within minikube

```shell
minikube ssh -- curl -Lo /tmp/kernel-headers-linux-4.19.94.tar.lz4 https://storage.googleapis.com/minikube-kernel-headers/kernel-headers-linux-4.19.94.tar.lz4

minikube ssh -- sudo mkdir -p /lib/modules/4.19.94/build

minikube ssh -- sudo tar -I lz4 -C /lib/modules/4.19.94/build -xvf /tmp/kernel-headers-linux-4.19.94.tar.lz4

minikube ssh -- rm /tmp/kernel-headers-linux-4.19.94.tar.lz4
```

c. build docker image within minikube

```shell
# change the destination to the docker in minikube
eval $(minikube docker-env)

# confirm docker context
docker context ls

# build images
docker build -t hello-ebpf-api:1.0.0 ./app-api
docker build -t bcc:1.0.0 ./bcc

# revert the destination to the local docker
eval $(minikube docker-env -u)

# confirm docker context
docker context ls
```

d. deploy to minikube

```shell
# create namespace
kubectl create namespace hello-ebpf

# apply manifest
kubectl apply -f k8s/api.yml --namespace=hello-ebpf

# get service url
minikube service hello-ebpf-api --url --namespace=hello-ebpf
```

e. run ebpf program

```shell
minikube ssh -- docker run --rm --privileged -v /lib/modules:/lib/modules:ro -v /usr/src:/usr/src:ro -v /etc/localtime:/etc/localtime:ro --workdir /usr/share/bcc/tools bcc:1.0.0 ./monitor_tcp_v4_connect.py
```

#### 1.2 Without eBPF programs

a. start minikube

```shell
minikube start --driver=virtualbox
```

b. build docker image in minikube

```shell
# change the destination to the docker in minikube
eval $(minikube docker-env)

# confirm docker context
docker context ls

# build image
docker build -t hello-ebpf-api:1.0.0 ./app-api

# revert the destination to the local docker
eval $(minikube docker-env -u)

# confirm docker context
docker context ls
```

c. deploy to minikube

```shell
# create namespace
kubectl create namespace hello-ebpf

# apply manifest
kubectl apply -f k8s/api.yml --namespace=hello-ebpf

# get service url
minikube service hello-ebpf-api --url --namespace=hello-ebpf
```

### 2. Stop the application

a. delete resources

```shell
# delete service
kubectl delete service hello-ebpf-api --namespace=hello-ebpf

# delete deployment
kubectl delete deployment hello-ebpf-api --namespace=hello-ebpf

# delete namespace
kubectl delete namespace hello-ebpf
```

b. stop minikube
```shell
minikube stop
```

c. delete minikube
```shell
minikube delete
```

## Contents

### Customs

- tools/monitor_tcp_v4_connect.py: Example of monitor tcp v4 connect.
- tools/ebpf_verifier_error.py: Example of an error in the ebpf program. <!-- TODO: 実装 -->

### Includes

See [iovisor/bcc](https://github.com/iovisor/bcc#contents).
