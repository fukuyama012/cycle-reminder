#!/usr/bin/env bash

go get -u golang.org/x/lint/golint

#golint ${GOPATH}/src/github.com/fukuyama012/cycle-reminder/service/web/app/controllers/
#golint ${GOPATH}/src/github.com/fukuyama012/cycle-reminder/service/web/app/models/
golint ${GOPATH}/src/github.com/fukuyama012/cycle-reminder/service/web/app/services/
