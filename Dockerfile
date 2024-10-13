#FROM mcr.microsoft.com/windows/servercore:ltsc2019
#
#WORKDIR /uploadserver
#
#COPY main.exe .
#COPY public ./public
#
#EXPOSE 3333
#
#CMD ["main.exe"]
# ENTRYPOINT ["./main.exe"]


FROM alpine:latest
WORKDIR /uploadserver
COPY . .
EXPOSE 3333
ENTRYPOINT ["./main"]
