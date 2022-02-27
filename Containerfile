FROM registry.access.redhat.com/ubi8/ubi:8.5 as builder
RUN curl -LO https://go.dev/dl/go1.17.7.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.7.linux-amd64.tar.gz
#RUN yum -y install golang
RUN mkdir -p /buildroot
COPY main.go /buildroot
COPY go.mod /buildroot
RUN mkdir /buildroot/alphavantage
COPY alphavantage/alphavantage.go /buildroot/alphavantage
RUN mkdir /buildroot/stocktickerresponse
COPY stocktickerresponse/stocktickerresponse.go /buildroot/stocktickerresponse
RUN cd /buildroot && CGO_ENABLED=0 /usr/local/go/bin/go build -o /stockticker main.go

FROM scratch
COPY --from=builder /stockticker /stockticker
ENTRYPOINT ["/stockticker"]