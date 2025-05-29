#!/bin/bash -xe

dnf -y module enable nodejs:12
dnf -y install nodejs
npm install markdownlint@0.26.0 markdownlint-cli2@0.4.0 --save-dev
