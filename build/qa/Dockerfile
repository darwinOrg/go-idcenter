#-alpine：指代镜像比较小
# FROM public.ecr.aws/bitnami/golang:1.19
FROM golang:1.19 as builder

LABEL stage=gobuilder

# 环境变量
ENV GO111MODULE=on \
    GOPRIVATE=dghire.com \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# GOPROXY=https://goproxy.cn,direct \

WORKDIR /application

#COPY go.mod , go.sum and download the dependencied
COPY . .
RUN go mod download

RUN go build -ldflags "-s -w" -o /application/build/qa-go-whale main.go

FROM ubuntu:22.04_stable
# FROM public.ecr.aws/docker/library/golang:1.19-alpine

WORKDIR /target

# 复制编译后的程序
COPY --from=builder /application/build/qa-go-whale /target/qa-go-whale
COPY --from=builder /application/conf/ /target/conf
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080
ENTRYPOINT ["/target/qa-go-whale"]
