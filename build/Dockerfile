FROM jetbrains/teamcity-minimal-agent as base
RUN apt-get update -q && \
    apt-get install -qy software-properties-common && \
    add-apt-repository ppa:longsleep/golang-backports && \
    apt-get update -q && \
    apt-get install -yq git golang && \
    apt-get upgrade -yq && \
    apt-get clean -q all

FROM base as build-buildtools
WORKDIR /src/app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go test -json ./...
RUN go build -o /bin/buildtools ./cmd/buildtools

FROM base as run
COPY ./build/run-docker.sh ./services
COPY --from=build-buildtools  /bin/buildtools /bin/buildtools
RUN apt-get install -qy apt-transport-https ca-certificates && \
    \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" && \
    add-apt-repository ppa:longsleep/golang-backports && \
    curl -sL https://deb.nodesource.com/setup_12.x | bash - && \
    \
    apt-cache policy docker-ce && \
    apt-get update -q && \
    apt-get install -yq  \
        docker-ce \
        docker-ce-cli \
        containerd.io \
        nodejs && \
    curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash && \
    \
    apt-get clean -q all && \
    \
    usermod -aG docker buildagent