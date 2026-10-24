# create build 
FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod ./
# COPY go.sum ./   # uncomment once the project has dependencies
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /reward-point-bin ./cmd/reward-point

# run build 
FROM scratch
COPY --from=build /reward-point-bin /reward-point
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/reward-point"]
