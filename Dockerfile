FROM alpine

ENV PORT 80
EXPOSE 80

RUN apk update && apk --no-cache add ca-certificates && update-ca-certificates

COPY build/helloaddon-linux /

ENTRYPOINT ["/helloaddon-linux", "-d"]
