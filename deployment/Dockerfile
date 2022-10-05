FROM golang:latest as builder

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY ./dist/docs /app
COPY ./dist/lbc /app

EXPOSE 80
EXPOSE 443

CMD ["/app/lbc"]