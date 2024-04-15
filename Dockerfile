FROM debian:bookworm-slim
ENV TZ=Asia/Shanghai
RUN apk update \
    && apk add tzdata \
    && echo "${TZ}" > /etc/timezone \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && rm /var/cache/apk/* \
ADD cLink /opt/cLink
ADD example.yaml /opt/config.yaml
WORKDIR /opt
CMD ["./cLink", "sync"]
