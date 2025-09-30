ARG IMG=gcr.io/distroless/static-debian11
FROM $IMG:nonroot

COPY example /usr/bin/local/example
COPY internal/db/_migrations /_migrations

ENTRYPOINT ["/usr/bin/local/example"]
