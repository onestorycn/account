FROM registry.cn-hangzhou.aliyuncs.com/meetwhale/golang:base


WORKDIR $GOPATH/src/whale-market

COPY ./ $GOPATH/src/whale-market
#RUN apt-get update && \
#    apt-get install -y --no-install-recommends \
#    vim && \
#    rm -rf /var/lib/apt/lists/* && \
#    apt-get clean
RUN go build ./grpc/main.go

EXPOSE 9994
ENTRYPOINT ["./main","--server_address=0.0.0.0:50051", "--registry=mdns"]
