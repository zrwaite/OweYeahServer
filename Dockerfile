FROM golang:1.18-alpine
WORKDIR /app
COPY ./ ./
RUN go mod download
RUN go build -o ./clueless-cartel-server

EXPOSE 8007

RUN cat .env

CMD [ "./clueless-cartel-server" ]
