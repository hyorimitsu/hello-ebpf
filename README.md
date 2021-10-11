Hello eBPF
---

This is a sample of running eBPF program in kubernetes.

## Description

This program is a sample to run the eBPF program (by BPF Compiler Collection) in minikube.
The configuration is as shown in the figure below.

![architecture](https://github.com/hyorimitsu/hello-ebpf/blob/master/doc/img/architecture.png)

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

- tools/[monitor_tcp_v4_connect.py](https://github.com/hyorimitsu/hello-ebpf/blob/master/bcc/tools/monitor_tcp_v4_connect.py): Example of monitor tcp v4 connect.

#### Output Example

```shell
PID    COMM         SADDR            DADDR            DPORT
5066   coredns      127.0.0.1        127.0.0.1        8080
4115   kubelet      127.0.0.1        127.0.0.1        10259
4115   kubelet      192.168.99.116   192.168.99.116   8443
5066   coredns      127.0.0.1        127.0.0.1        8080
4126   kubelet      172.17.0.1       172.17.0.2       8181
4115   kubelet      172.17.0.1       172.17.0.2       8080
4114   <...>        192.168.99.116   192.168.99.116   8443
4602   coredns      127.0.0.1        127.0.0.1        8080
4126   kubelet      192.168.99.116   192.168.99.116   8443
4597   <...>        127.0.0.1        127.0.0.1        8080
3580   kube-apiserv 127.0.0.1        127.0.0.1        2379
...
# API is connecting to the https://github.com/ when accessing the http://${Service IP}:${Service Port}/?url=https://github.com/ in the browser.
7975   app-api      172.17.0.4       13.114.40.48     443
...
```

#### Output Description

|Name|Description|
|----|-----------|
|PID|Process ID|
|COMM|Exec Command|
|SADDR|Sender IP|
|DADDR|Destination IP|
|DPORT|Destination Port|

#### Output Confirm

```shell
$ host github.com
github.com has address 13.114.40.48
...

$ kubectl get pods --output=wide --namespace=hello-ebpf
NAME                             READY   STATUS    RESTARTS   AGE   IP           NODE       NOMINATED NODE   READINESS GATES
hello-ebpf-api-5fd6fc94c-7s4vc   1/1     Running   0          14m   172.17.0.3   minikube   <none>           <none>
hello-ebpf-api-5fd6fc94c-bl2jn   1/1     Running   0          14m   172.17.0.5   minikube   <none>           <none>
hello-ebpf-api-5fd6fc94c-cl5tw   1/1     Running   0          14m   172.17.0.4   minikube   <none>           <none>
```

- tools/[ebpf_verifier_error.py](https://github.com/hyorimitsu/hello-ebpf/blob/master/bcc/tools/ebpf_verifier_error.py): Example of an error in the ebpf program.

#### Output Example

```shell
bpf: Failed to load program: Permission denied
0: (b7) r1 = 4660
1: (b7) r2 = 43981
2: (63) *(u32 *)(r1 +0) = r2
R1 invalid mem access 'inv'

HINT: The invalid mem access 'inv' error can happen if you try to dereference memory without first using bpf_probe_read() to copy it to the BPF stack. Sometimes the bpf_probe_read is automatic by the bcc rewriter, other times you'll need to be explicit.

Traceback (most recent call last):
  File "/usr/share/bcc/tools/ebpf_verifier_error.py", line 17, in <module>
    fn_name="execve_hook")
  File "/usr/lib/python2.7/dist-packages/bcc/__init__.py", line 650, in attach_kprobe
    fn = self.load_func(fn_name, BPF.KPROBE)
  File "/usr/lib/python2.7/dist-packages/bcc/__init__.py", line 391, in load_func
    (func_name, errstr))
Exception: Failed to load BPF program execve_hook: Permission denied
ssh: Process exited with status 1
```

### Includes

See [iovisor/bcc](https://github.com/iovisor/bcc#contents).
