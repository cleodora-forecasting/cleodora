FROM scratch
COPY cleosrv /
COPY cleosrv.example.yml /data/cleosrv.yml
EXPOSE 8080/tcp
VOLUME /data
ENTRYPOINT [ \
    "/cleosrv", \
    "--address", \
    "0.0.0.0:8080", \
    "--database", \
    "/data/cleosrv.db", \
    "--config", \
    "/data/cleosrv.yml" \
    ]
CMD []
