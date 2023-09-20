FROM registry.cn-shanghai.aliyuncs.com/star_base/golang-x86:1.19.1 as builder

LABEL stage=gobuilder

# 环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux    \
    GOARCH=arm64

WORKDIR /application

COPY . .
RUN go env
RUN go mod download

RUN go build -ldflags "-s -w" -o /application/build/qa-go-idcenter main.go

FROM registry.cn-shanghai.aliyuncs.com/star_base/ubuntu22:22.04_stable

WORKDIR /target

# 复制编译后的程序
COPY --from=builder /application/build/qa-go-idcenter /target/qa-go-idcenter
COPY --from=builder /application/conf/ /target/conf
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080
ENTRYPOINT ["/target/qa-go-idcenter"]