FROM golang:latest

#create and fill app folder
WORKDIR /app
COPY go.mod go.sum ./
COPY *.go ./
COPY todo/*.go ./todo/
COPY config/*.go ./config/
COPY config/*.yaml ./config/
COPY fake/*.go ./fake/

#dl libraries
RUN go mod download

RUN go build

#documenting ports to forward
EXPOSE 8080 80

#run the app
ENTRYPOINT ["go", "run", "."]