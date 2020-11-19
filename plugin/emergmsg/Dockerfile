FROM golang:alpine
RUN apk add git ca-certificates
RUN git clone https://github.com/coredns/coredns /coredns
WORKDIR /coredns/plugin/emergmsg
COPY . .
WORKDIR /coredns
RUN echo "emergmsg:emergmsg" >> plugin.cfg; \
    go generate && go build -v


FROM scratch

COPY --from=0 /etc/ssl/certs /etc/ssl/certs
COPY --from=0 /coredns/coredns /coredns

EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]