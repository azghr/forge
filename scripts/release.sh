#!/bin/bash

PKG=$1
VERSION=$2

git tag ${PKG}/v${VERSION}
git push origin ${PKG}/v${VERSION}