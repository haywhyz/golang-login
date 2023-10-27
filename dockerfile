# FROM golang:1.19-alpine

# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./

# RUN go mod download

# COPY *.go ./

# COPY backend.env ./.env

# RUN go test
# RUN go build -o /main

# EXPOSE 4000

# CMD [ "/main" ]



FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

# Run tests and build and run
RUN go test -v

RUN go build -o /main

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /main /main

EXPOSE 4000

USER nonroot:nonroot

CMD [ "/main" ]