FROM alpine:3.12
COPY ./$path/sfdata /app
RUN chmod 555 /app/sfdata
WORKDIR /app
EXPOSE 3000
CMD ["/app/sfdata"]
