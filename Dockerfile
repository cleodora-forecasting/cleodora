FROM scratch
COPY dist/cleosrv_linux_amd64_v1/cleosrv /
EXPOSE 8080/tcp
ENTRYPOINT ["/cleosrv", "--address", "0.0.0.0:8080"]
