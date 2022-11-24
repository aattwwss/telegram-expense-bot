FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /telegram-expense-bot

#EXPOSE 8080

CMD [ "/telegram-expense-bot" ]