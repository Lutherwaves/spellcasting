FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/SERVICENAME .

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/SERVICENAME /SERVICENAME
USER nonroot:nonroot
EXPOSE SERVICEPORT
ENTRYPOINT ["/SERVICENAME", "server"]
