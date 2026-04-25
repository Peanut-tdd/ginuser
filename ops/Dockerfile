# 构建命令：在项目根目录下执行 docker build -f ops/Dockerfile -t ginuser .
FROM golang:1.25-alpine as builder

# 设置代理和环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 1. 先拷贝依赖文件，利用 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download

# 2. 拷贝源代码
COPY . .

# 3. 编译（输出到当前目录下的 main）
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# 4. 准备发布目录
RUN mkdir publish \
    && cp main publish/ \
    && cp -r etc publish/

# 运行阶段
FROM alpine:3.19

# 设置时区（可选，但通常需要）
RUN apk --no-cache add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 从 build 阶段拷贝文件
COPY --from=builder /app/publish .

# 指定运行时环境变量
ENV GIN_MODE=release
EXPOSE 5001

ENTRYPOINT ["./main"]






