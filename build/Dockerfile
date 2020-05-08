FROM jetbrains/teamcity-agent as base
RUN add-apt-repository ppa:longsleep/golang-backports && \
    curl -sL https://deb.nodesource.com/setup_12.x | bash - && \
    apt-get update -q && \
    apt-get install -yq golang nodejs && \
    apt-get upgrade -yq && \
    apt-get clean -q all && \
    curl -sL https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash

FROM base as build-buildtools
WORKDIR /src/app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go test -json ./...
RUN go build -o /bin/buildtools ./cmd/buildtools

FROM base as run
COPY --from=build-buildtools  /bin/buildtools /bin/buildtools