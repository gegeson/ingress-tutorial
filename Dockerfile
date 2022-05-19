FROM golang:1.18 AS build

ARG LDFLAGS="-w"
ENV CGO_ENABLED=0

WORKDIR /go/app
COPY . /go/app

RUN go build -ldflags="${LDFLAGS}" -o /usr/local/bin/sample .

FROM alpine:3.15
RUN apk --no-cache add ca-certificates tzdata
COPY --from=build /usr/local/bin /usr/local/bin

ENTRYPOINT ["sample"]