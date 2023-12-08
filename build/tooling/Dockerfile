FROM golang:1.21-bullseye@sha256:2148dd89f087c3e0efe8f1337ef325fede10bb6f9109cc367b9b601fb3299379

RUN GO111MODULE=on go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.9.0
RUN GO111MODULE=on go install k8s.io/code-generator/cmd/conversion-gen@v0.25.4

RUN git config --global --add safe.directory "*"

WORKDIR /eraser
