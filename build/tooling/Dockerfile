FROM golang:1.20-bullseye@sha256:c8af7d1086afa15ba41f257efc6382f42e29f9e1401e64fc69ab3d076f8a6204

RUN GO111MODULE=on go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.9.0
RUN GO111MODULE=on go install k8s.io/code-generator/cmd/conversion-gen@v0.25.4

RUN git config --global --add safe.directory "*"

WORKDIR /eraser
