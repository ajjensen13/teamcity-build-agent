ARG GO_VERSION=1
ARG ALPINE_VERSION=3.11
ARG TEAMCITY_VERSION=2020.2.1

FROM jetbrains/teamcity-agent:${TEAMCITY_VERSION}-linux-sudo as base
USER root

RUN add-apt-repository ppa:longsleep/golang-backports && \
    curl -sL https://deb.nodesource.com/setup_12.x | bash - && \
    curl -sL https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    curl -sL https://storage.googleapis.com/download.dartlang.org/linux/debian/dart_stable.list > /etc/apt/sources.list.d/dart_stable.list && \
    curl -sL https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash && \
    apt-get update && \
    apt-get install -yq golang nodejs dart gcc g++ make && \
    apt-get clean -q all
ENV PATH="/root/.pub-cache/bin:/usr/lib/dart/bin:$PATH"
RUN pub global activate webdev
RUN npm install -g @angular/cli

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as build-buildtools
ARG GCFLAGS="-c 1"
ENV CGO_ENABLED=0
ENV GO111MODULE=on
WORKDIR /src/app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go test -json ./...
RUN GOOS=linux GOARCH=amd64 go build -gcflags "$GCFLAGS" -o /bin/buildtools ./cmd/buildtools

FROM base as run
COPY --from=build-buildtools  /bin/buildtools /bin/buildtools