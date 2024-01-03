FROM golang:1.18 AS builder

# Set the working directory to golang working space
WORKDIR /riotpot

# Copy the dependencies into the image
COPY go.mod ./
COPY go.sum ./

# download all the dependencies
RUN go mod download

# Copy the files from the previous build into the ui folder
COPY ui ui/

# Copy everything into the image
# Copy only the app files in the image
COPY api api/
COPY cmd cmd/
COPY pkg pkg/
COPY plugin plugin/

# Copy static files
COPY statik/ statik/
COPY build build/

ADD Makefile .
RUN make compile

FROM golang:1.18 AS release

#ENV GIN_MODE=release

WORKDIR /riotpot

# Copy the dependencies into the image
COPY --from=builder /riotpot/bin/riotpot ./

# UI port
EXPOSE 3000

CMD ["./riotpot"]