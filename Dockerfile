FROM golang:alpine as builder
RUN apk update && \
    apk upgrade && \
    apk add git
RUN apk add build-base
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go env
RUN go build -o main .
FROM alpine
COPY --from=builder /build/main /app/
WORKDIR /app
RUN ls -a
CMD ["./main"]