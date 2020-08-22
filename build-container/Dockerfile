FROM alpine
ARG dbmongoDir
COPY ./$dbmongoDir/dbmongo /app/dbmongo
RUN chmod 555 /app/dbmongo
WORKDIR /app
EXPOSE 3000
CMD ["/app/dbmongo"]
