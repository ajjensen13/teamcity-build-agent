# template #

## Overview ##
This is a template project used for containerized 
Golang projects that are deployed on Kubernetes using Helm.

## Structure ##
The code loosely follows the structure prescribed by [github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout).

### Go Directories

#### `/cmd`

Main applications for this project.

The directory name for each application should match the name of the executable you want to have (e.g., `/cmd/myapp`).

#### `/internal`

Private application and library code. 

#### `/pkg`

Library code that's ok to use by external applications (e.g., `/pkg/mypubliclib`).

### Web Application Directories

#### `/web`

Web application specific components: static web assets, server side templates and SPAs.

### Common Application Directories

#### `/scripts`

Scripts to perform various build, install, analysis, etc operations.

#### `/build`

Packaging and Continuous Integration.

Put your cloud (AMI), container (Docker), OS (deb, rpm, pkg) package configurations and scripts in the `/build/package` directory.

#### `/deploy`

IaaS, PaaS, system and container orchestration deployment configurations and templates (docker-compose, kubernetes/helm, mesos, terraform, bosh).