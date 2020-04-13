FROM golang:1.13.10-buster as backend-builder

ARG CI
ARG DRONE
ARG DRONE_TAG
ARG DRONE_COMMIT
ARG DRONE_BRANCH
ARG DRONE_PULL_REQUEST

ARG SKIP_BACKEND_TEST
ARG BACKEND_TEST_TIMEOUT

ADD backend /build/backend
ADD .git/ /build/backend/.git/
WORKDIR /build/backend

ENV GOFLAGS="-mod=vendor"

# run tests
RUN \
    cd app && \
    if [ -z "$SKIP_BACKEND_TEST" ] ; then \
        go test -p 1 -timeout="${BACKEND_TEST_TIMEOUT:-300s}" -covermode=count -coverprofile=/profile.cov_tmp ./... && \
        cat /profile.cov_tmp | grep -v "_mock.go" > /profile.cov ; \
        golangci-lint run --config ../.golangci.yml ./... ; \
    else echo "skip backend tests and linter" ; fi

# if DRONE presented use DRONE_* git env to make version
RUN \
    if [ -z "$DRONE" ] ; then echo "runs outside of drone" && version="$(git rev-parse --short HEAD)" ; \
    else version=${DRONE_TAG}${DRONE_BRANCH}${DRONE_PULL_REQUEST}-${DRONE_COMMIT:0:7}-$(date +%Y%m%d-%H:%M:%S) ; fi && \
    echo "version=$version" && \
    go build -o remark42 -ldflags "-X main.revision=${version} -s -w" ./app

FROM node:10.11-alpine as frontend-builder-deps

ARG CI
ENV HUSKY_SKIP_INSTALL=true

RUN apk add --no-cache --update git
ADD frontend/package.json /srv/frontend/package.json
ADD frontend/package-lock.json /srv/frontend/package-lock.json
RUN cd /srv/frontend && CI=true npm ci

FROM node:10.11-alpine as frontend-builder

ARG CI
ARG SKIP_FRONTEND_TEST
ARG NODE_ENV=production

COPY --from=frontend-builder-deps /srv/frontend/node_modules /srv/frontend/node_modules
ADD frontend /srv/frontend
RUN cd /srv/frontend && \
    if [ -z "$SKIP_FRONTEND_TEST" ] ; then npx run-p check lint lint:style test build ; \
    else echo "skip frontend tests and lint" ; npm run build ; fi && \
    rm -rf ./node_modules

# merge the build
FROM 80x86/base-debian:buster-slim-amd64 as stage

RUN mkdir /stage

WORKDIR /stage

COPY docker-init.sh ./entrypoint.sh
COPY backend/scripts/backup.sh ./usr/local/bin/backup
COPY backend/scripts/restore.sh ./usr/local/bin/restore
COPY backend/scripts/import.sh ./usr/local/bin/import

RUN chmod +x ./entrypoint.sh ./usr/local/bin/backup ./usr/local/bin/restore ./usr/local/bin/import

COPY --from=backend-builder /build/backend/remark42 ./srv/remark42
COPY --from=frontend-builder /srv/frontend/public/ ./srv/web

COPY docker-init.sh ./srv/init.sh

RUN mkdir ./srv/var
RUN chown -R app:app ./srv
RUN mkdir -p ./usr/bin && ln -s ./srv/remark42 ./usr/bin/remark42
RUN chmod +x ./srv/init.sh

# final step
FROM 80x86/base-debian:buster-slim-amd64 as final

COPY --from=stage /stage/ /

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD curl --fail http://localhost:8080/ping || exit 1

WORKDIR /srv

CMD ["/srv/remark42", "server"]
