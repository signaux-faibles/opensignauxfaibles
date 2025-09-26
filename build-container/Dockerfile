FROM alpine:3.19
COPY --chmod=555 sfdata /app/sfdata
COPY local_data /local_data
WORKDIR /app
EXPOSE 3000
CMD ["/app/sfdata"]
