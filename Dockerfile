FROM golang:alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apk update && apk add bash
RUN apk add postgresql-client

# wait postgres
RUN chmod +x wait-for-postgres.sh
# RUN dos2unix wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o subscription_microservice ./cmd/main.go

CMD ["./subscription_microservice"]