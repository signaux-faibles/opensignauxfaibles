FROM alpine:3.12
COPY ./$path/dbmongo /app
RUN chmod 555 /app/dbmongo
WORKDIR /app
EXPOSE 3000
CMD ["/app/dbmongo"]
