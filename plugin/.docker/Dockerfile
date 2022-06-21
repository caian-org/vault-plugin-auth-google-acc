FROM golang:1.18.1-alpine3.15 AS go

ARG PLUGIN_PATH=vault-plugin-auth-google-acc
ENV PLUGIN_PATH $PLUGIN_PATH

RUN apk update \
    && apk add --no-cache make \
    && rm -rf /var/cache/apk/*


FROM go AS plugin-dependencies
WORKDIR /go/plugin
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go list -m all \
    | tail -n +2 \
    | cut -f 1 -d " " \
    | awk 'NF{print $0 "/..."}' \
    | CGO_ENABLED=0 GOOS=linux xargs -n1 \
    go build -v -installsuffix cgo -i; echo "done"


FROM plugin-dependencies AS plugin-build
WORKDIR /go/plugin
COPY . .
RUN make container-binary


FROM vault:1.9.7 as vault
COPY --from=plugin-build "/go/plugin/bin/${PLUGIN_PATH}" /etc/vault/plugins
RUN chown vault:vault "/etc/vault/plugins/${PLUGIN_PATH}"
