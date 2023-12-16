FROM node:19.4.0-alpine AS ui

# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# add app
COPY ui/ .
RUN npm install --silent

# Creates the build
RUN npm run build --omit=dev

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
COPY plugins plugins/

# Copy static files
COPY statik/ statik/
COPY build build/

ADD Makefile .
RUN make

FROM gcr.io/distroless/base-debian10 AS release

#ENV GIN_MODE=release

WORKDIR /riotpot

# Copy the dependencies into the image
COPY --from=builder /riotpot/bin/riotpot ./

# UI port
EXPOSE 2022

CMD ["./riotpot"]