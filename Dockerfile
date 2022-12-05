FROM golang:1.19 as builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /server

FROM gcr.io/distroless/base-debian10
WORKDIR /usr/src/app
COPY --from=builder /server .
CMD ["/usr/src/app/server"]
