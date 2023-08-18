#!/bin/bash

NUM=$1
path=$2

cat <<EOF >"$path"/daemonset.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: preload
spec:
  selector:
    matchLabels:
      name: preload
  template:
    metadata:
      labels:
        name: preload
    spec:
      initContainers:
EOF

for ((i=1; i<=NUM; i++))
do
cat <<EOF >>"$path"/daemonset.yaml
        - name: preload-$i
          image: gcr.io/k8s-testimages/perf-tests-util/containerd:v0.0.1
          command: ["sh", "-c", "ctr -n=k8s.io image pull asmehrotra.azurecr.io/vulnimage:v$i 2>&1 | tee /var/log/cl2-image-preload-$i.log"]
          volumeMounts:
          - name: containerd
            mountPath: /run/containerd
          - name: logs-volume
            mountPath: /var/log
EOF
done

cat <<EOF >>"$path"/daemonset.yaml
      volumes:
        - name: containerd
          hostPath:
            path: /run/containerd
        - name: logs-volume
          hostPath:
            path: /var/log
      containers:
        - name: pause
          image: gcr.io/google_containers/pause
EOF
