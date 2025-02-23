FROM golang:1.24-alpine AS builder

RUN apk add --no-cache nodejs npm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN rm -rf client/dist
RUN rm -rf bin
RUN mkdir -p bin

RUN go mod download 
RUN go build -o bin/ .
RUN cd client && npm install && npm run build
RUN mkdir -p bin/client/dist 
RUN cp -r client/dist/* bin/client/dist

FROM alpine:latest

RUN apk add --no-cache --purge -uU libreoffice \
    && rm -rf /var/cache/apk/* /tmp/*

COPY ./fonts/* /usr/local/share/fonts/

RUN mkdir -p /app
WORKDIR /app

COPY --from=builder /app/bin /app/
COPY ./templates/* /app/templates/

EXPOSE 80

ENTRYPOINT [ "./anrechnungsstundenberechner" ]
CMD ["serve", "--http", "0.0.0.0:80"]