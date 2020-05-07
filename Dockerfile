FROM golang:1.14 AS build
WORKDIR /go/src
COPY pkg/server ./pkg/server
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o identityserver .

FROM scratch AS runtime
COPY --from=build /go/src/identityserver ./
EXPOSE 8080/tcp
ENTRYPOINT ["./identityserver"]
