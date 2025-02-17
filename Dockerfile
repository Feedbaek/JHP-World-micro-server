FROM golang:1.22-alpine

# go build
WORKDIR /app
COPY . .

# install packages
RUN apk add --no-cache g++ util-linux

# go mod
RUN go mod tidy && go build -o go_is_awesome

# expose port
EXPOSE 8080

# run
ENTRYPOINT ["./go_is_awesome"]
