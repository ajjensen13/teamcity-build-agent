ARG GO_VERSION=1
ARG ALPINE_VERSION=3.11

FROM jetbrains/teamcity-agent as base
USER root
# RUN echo "root ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers.d/root_nopasswd && chmod 0440 /etc/sudoers.d/root_nopasswd

RUN apt-get update -q && \
    apt-get upgrade -yq
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
COPY ./build/setup-config.sh /services/
COPY ./build/run-services.sh /run-services.sh