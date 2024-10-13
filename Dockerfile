FROM mcr.microsoft.com/windows/servercore:ltsc2022

WORKDIR /uploadserver

COPY main.exe .
COPY public ./public

EXPOSE 3333

CMD ["main.exe"]
# ENTRYPOINT ["./main.exe"]
