FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD apiserver apiserver
EXPOSE 80
ENTRYPOINT ["/apiserver"]
