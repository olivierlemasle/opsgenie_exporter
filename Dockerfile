# Step 1: build executable binary
FROM golang:1.18-alpine as builder

ENV CGO_ENABLED=0
WORKDIR /go/src/app

# Copy only go.mod and go.sum files before downloading modules,
# so that modules can be cached if these two files do not change.
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY pkg pkg

ARG VERSION
ARG SOURCE_COMMIT
ARG SOURCE_BRANCH

RUN go install \
    -ldflags=" \
    -X github.com/prometheus/common/version.Version=${VERSION} \
    -X github.com/prometheus/common/version.Revision=${SOURCE_COMMIT} \
    -X github.com/prometheus/common/version.Branch=${SOURCE_BRANCH} \
    -X github.com/prometheus/common/version.BuildUser=$(whoami) \
    -X github.com/prometheus/common/version.BuildDate=$(date +%Y%m%d-%T)" \
    .

# Step 2: build image
FROM gcr.io/distroless/static:latest

COPY --from=builder /go/bin/opsgenie_exporter /

USER nobody

EXPOSE 3000

CMD ["/opsgenie_exporter"]
