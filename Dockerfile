FROM golang:1.19.4
WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o executable.run
RUN chmod +x executable.run

CMD [ "./executable.run" ]