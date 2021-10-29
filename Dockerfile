FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

LABEL maintainer="lzuccarelli@tfd.ie"

RUN mkdir -p /go/src /go/bin && chmod -R 0755 /go
COPY build/microservice uid_entrypoint.sh /go/ 
WORKDIR /go

USER 1001

ENTRYPOINT [ "./uid_entrypoint.sh" ]

# This will change depending on each microservice entry point
CMD ["./microservice"]
