# Build the njmon_exporter binary
FROM golang:1.19 as builder

WORKDIR /workspace

# Copy the go source
COPY go.mod go.sum collector.go lastseen.go listener.go main.go .
ADD config ./config

# Introduce the build arg check in the end of the build stage
# to avoid messing with cached layers
ARG VERSION

# Fetch modules via the proxy
ENV GOPROXY=http://plonexus01.westernpower.co.uk:8081/repository/go-proxy/
ENV GOSUMDB="sum.golang.org http://plonexus01.westernpower.co.uk:8081/repository/go-sum-proxy/"
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -ldflags "-X main.buildVersion=${VERSION} -X main.buildDate=`date -u +%Y-%m-%d`" -o njmon_exporter .

RUN test -n "$VERSION" || (echo "VERSION not set" && false)

# Use the scratch image since we only need the go binary
FROM scratch

ARG VERSION

LABEL name=njmon_exporter \
      vendor='National Grid Electricity Distribution' \
      version=$VERSION \
      release=$VERSION \
      description='njmon exporter image' \
      summary='Export metrics from njmon on AIX'

ENV USER_ID=1001

WORKDIR /
COPY --from=builder /workspace/njmon_exporter .
COPY config/njmon_exporter.yml .
USER ${USER_ID}

EXPOSE 9772

CMD ["/njmon_exporter --config njmon_exporter.yml"]
