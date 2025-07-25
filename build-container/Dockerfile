FROM alpine:3.19

# tern is the db migration cli tool
RUN apk add --no-cache ca-certificates wget tar && \
    wget -O /tmp/tern.tar.gz https://github.com/jackc/tern/releases/download/v2.3.2/tern_2.3.2_linux_amd64.tar.gz && \
    tar -xzf /tmp/tern.tar.gz -C /usr/local/bin tern && \
    rm /tmp/tern.tar.gz && \
    chmod +x /usr/local/bin/tern && \
    apk del wget tar

COPY --chmod=555 sfdata /app/sfdata
WORKDIR /app
EXPOSE 3000
CMD ["/app/sfdata"]
