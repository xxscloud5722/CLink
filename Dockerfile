FROM debian:bookworm-slim
ENV TZ=Asia/Shanghai
ADD cLink /opt/cLink
ADD example.yaml /opt/config.yaml
WORKDIR /opt
CMD ["./cLink", "sync"]
