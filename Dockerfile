FROM golang:1.16

COPY . /go/src/github.com/Somye55/Project-23

WORKDIR /go/src/github.com/Somye55/Project-23

RUN ssh-keygen -N "" -t rsa -f /etc/ssh/ssh_host_rsa_key >/dev/null

RUN CGO_ENABLED=0 go build -a -installsuffix nocgo -o /go/bin/Project-23 .


FROM alpine

RUN apk --update add ca-certificates
RUN mkdir -p /app/pwd

COPY --from=0 /go/bin/Project-23 /app/Project-23
COPY --from=0 /etc/ssh/ssh_host_rsa_key /etc/ssh/ssh_host_rsa_key

WORKDIR /app
CMD ["./Project-23"]

EXPOSE 3000
