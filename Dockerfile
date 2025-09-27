ARG IMG=gcr.io/distroless/static-debian11
FROM $IMG:nonroot

COPY example /usr/bin/local/example

ENTRYPOINT ["/usr/bin/local/example"]
