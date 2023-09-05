FROM golang:1.21

WORKDIR /app

RUN apt install make

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN make build/service
RUN make build/client

CMD ["make", "run/service"]