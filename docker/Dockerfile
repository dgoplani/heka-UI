# build the frontend UI
FROM node:16.16.0-alpine AS frontend-builder
LABEL stage=frontend
WORKDIR /angular
COPY frontend .
RUN npm install
RUN npm run build-bloxconnect

# build the backend server binary
FROM golang:1.16-buster AS backend-builder
LABEL stage=backend
WORKDIR /go/src/github.com/Infoblox-CTO/heka-ui
COPY backend .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o bin/heka-ui ./cmd

# Final Container Image
FROM alpine:3.16 AS runner
LABEL com.infoblox.maintainer="infoblox"
LABEL maintainer="Infoblox <support@infoblox.com>"
ARG GIT_TAG
LABEL git_tag=$GIT_TAG
ARG BUILD_DATE
LABEL build_date=$BUILD_DATE

WORKDIR /bin

RUN apk update && apk upgrade
RUN apk add bash
RUN apk add nginx

COPY config/nginx.conf /etc/nginx/nginx.conf
COPY --from=frontend-builder /angular/dist/heka-ui /usr/local/angular
COPY --from=backend-builder /go/src/github.com/Infoblox-CTO/heka-ui/scripts/* ./
COPY --from=backend-builder /go/src/github.com/Infoblox-CTO/heka-ui/bin/heka-ui .

EXPOSE 26750
EXPOSE 26751

ENTRYPOINT ["heka-ui"]
