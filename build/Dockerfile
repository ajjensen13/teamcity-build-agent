FROM jetbrains/teamcity-agent as base
RUN curl -sL https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
RUN apt-get update -q && \
    apt-get upgrade -yq
RUN add-apt-repository ppa:longsleep/golang-backports && \
    apt-get update && \
    apt-get install -yq golang && \
    apt-get clean -q all

FROM base as build-buildtools
WORKDIR /src/app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go test -json ./...
RUN go build -o /bin/buildtools ./cmd/buildtools

FROM base as run
RUN curl -sL https://deb.nodesource.com/setup_12.x | bash - && \
    curl -sL https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    curl -sL https://storage.googleapis.com/download.dartlang.org/linux/debian/dart_stable.list > /etc/apt/sources.list.d/dart_stable.list && \
    apt-get update
RUN apt-get install -yq nodejs dart gcc g++ make man && \
    apt-get clean -q all
ENV PATH="$HOME/.pub-cache/bin:/usr/lib/dart/bin:$PATH"
RUN pub global activate webdev
COPY --from=build-buildtools  /bin/buildtools /bin/buildtools
COPY ./build/setup-config.sh /services/