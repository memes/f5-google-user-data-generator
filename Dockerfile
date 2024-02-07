FROM alpine:3.19.1 as ca
RUN apk --no-cache add ca-certificates-bundle=20191127-r5

FROM scratch
COPY --from=ca /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
COPY f5-google-user-data-generator /f5-google-user-data-generator
ENTRYPOINT ["/f5-google-user-data-generator"]
