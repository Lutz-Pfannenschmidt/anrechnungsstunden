FROM golang:1.24-alpine

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
RUN cp temlate.xlsx bin/ 
RUN cd client && npm install && npm run build
RUN mkdir -p bin/client/dist 
RUN cp -r client/dist/* bin/client/dist

EXPOSE 8090

# Run
CMD ["./bin/anrechnungsstundenberechner", "serve", "--http", "0.0.0.0:80"]