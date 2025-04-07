FROM golang:1.24-alpine


WORKDIR /app
RUN mkdir -p code
COPY dataCollector/ ./code/


WORKDIR /app/code
RUN go mod download

RUN go build -o /app/main ./cmd/.
WORKDIR /app

CMD ["sh", "-c", "/app/main /app/code/configs/app.ini"]