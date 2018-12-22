FROM golang:alpine as builder
RUN apk update && \
    apk upgrade && \
    apk add git
RUN apk add build-base
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN export GO111MODULE=on
RUN go build -mod=vendor -o main .
FROM alpine
COPY --from=builder /build/main /app/
ADD /google-key/server.crt /etc/ssl/certs/
# ADD /google-key/key.json /google-key/ conflict with secret Mountvolume Error
WORKDIR /app
RUN apk --update add ca-certificates
CMD ["./main"]