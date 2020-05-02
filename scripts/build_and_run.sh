#!/bin/bash

(
    go get github.com/Project-Wartemis/pw-backend/cmd/backend &&
    $GOPATH/bin/backend
)
