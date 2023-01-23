FROM alpine:3.12
COPY sfdata_broken /app
RUN chmod 555 /app/sfdata
WORKDIR /app
EXPOSE 3000
CMD ["/app/sfdata"]
