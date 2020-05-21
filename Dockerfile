FROM golang:1.13 as builder

RUN cd ..
RUN mkdir gitpr
WORKDIR gitpr
COPY . ./

RUN sed -i s/{serviceMode}/http/g configuration.yaml
RUN sed -i s/{servicePort}/9999/g configuration.yaml

ARG version=0.0.1
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags "-X main.version=$version" -o gitpr ./cmd/gitpr/main.go

FROM scratch

COPY --from=builder /go/gitpr/gitpr .
COPY --from=builder /go/gitpr/configuration.yaml configuration.yaml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./gitpr"]
