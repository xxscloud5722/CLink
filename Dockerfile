FROM debian:bookworm-slim
ADD cLink /opt/cLink
ADD example.yaml /opt/config.yaml
WORKDIR /opt
CMD ["clink", "sync"]
