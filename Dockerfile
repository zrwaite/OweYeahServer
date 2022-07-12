FROM golang:1.18-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./owemate-server

EXPOSE 8007

CMD [ "./owemate-server" ]
