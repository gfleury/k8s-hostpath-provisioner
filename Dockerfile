FROM golang:1.9.0 as builder
COPY . /go/src/github.com/gfleury/k8s-hostpath-provisioner/
WORKDIR /go/src/github.com/gfleury/k8s-hostpath-provisioner/
RUN make hostpath-provisioner

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/gfleury/k8s-hostpath-provisioner/hostpath-provisioner .

CMD ["./hostpath-provisioner"]
