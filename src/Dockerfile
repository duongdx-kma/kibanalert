# Stage 1: build artifacts stage
FROM public.ecr.aws/docker/library/golang:1.19.13 as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

# build artifacts
RUN CGO_ENABLED=0 GOOS=linux go build -o artifacts run.go

# Stage 2: Create a lightweight image for runtime
FROM public.ecr.aws/docker/library/alpine:3.18.5 as deploy
WORKDIR /app

# copy artifacts from build stage
COPY --from=builder /app/artifacts /app/run

# make artifacts executable
RUN chmod +x /app/run

# config container timezone
RUN apk --no-cache add tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
    echo "Asia/Ho_Chi_Minh" > /etc/timezone

EXPOSE 8080
CMD ["/app/run"]