FROM golang:alpine as builder
RUN apk update && \
    apk upgrade && \
    apk add git
RUN apk add build-base
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go env && export GO111MODULE=on
RUN go build -mod=vendor -v -o main .
FROM alpine
COPY --from=builder /build/main /app/
WORKDIR /app
VOLUME "/google-key"
CMD ["./main"]