FROM alpine:latest
WORKDIR /uploadserver
COPY . .
EXPOSE 3333
ENTRYPOINT ["./main"]
