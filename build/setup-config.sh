#!/usr/bin/env bash

mkdir -p "/data/teamcity_agent/conf"
PROPERTIES_FILE="/data/teamcity_agent/conf/buildAgent.properties"

command -v dart > /dev/null && \
  printf "dart.path=%s\n" "$(command -v dart)" >> "$PROPERTIES_FILE" && \
  printf "dart.version=%s\n" "$(dart --version 2>&1)" >> "$PROPERTIES_FILE"

command -v pub > /dev/null && \
  printf "pub.path=%s\n" "$(command -v pub)" >> "$PROPERTIES_FILE" && \
  printf "pub.version=%s\n" "$(pub version)" >> "$PROPERTIES_FILE"

command -v webdev > /dev/null && \
  printf "webdev.path=%s\n" "$(command -v webdev)" >> "$PROPERTIES_FILE" && \
  printf "webdev.version=%s\n" "$(webdev --version)" >> "$PROPERTIES_FILE"

command -v node > /dev/null && \
  printf "node.path=%s\n" "$(command -v node)" >> "$PROPERTIES_FILE" && \
  printf "node.version=%s\n" "$(node --version)" >> "$PROPERTIES_FILE"

command -v npm > /dev/null && \
  printf "npm.path=%s\n" "$(command -v npm)" >> "$PROPERTIES_FILE" && \
  printf "npm.version=%s\n" "$(npm --version)" >> "$PROPERTIES_FILE"

command -v go > /dev/null && \
  printf "go.path=%s\n" "$(command -v go)" >> "$PROPERTIES_FILE" && \
  printf "go.version=%s\n" "$(go version)" >> "$PROPERTIES_FILE"

command -v ng > /dev/null && \
  printf "ng.path=%s\n" "$(command -v ng)" >> "$PROPERTIES_FILE" && \
  printf "ng.version=%s\n" "$(ng version | grep "Angular CLI")" >> "$PROPERTIES_FILE"

