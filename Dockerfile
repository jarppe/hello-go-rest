#
# Build application:
#


FROM golang:1.13-alpine AS build


WORKDIR /build

# Dependencies:
COPY ./go.mod /build
COPY ./go.sum /build
RUN cd /build && go mod download


# Application sources:
COPY ./*.go /build

# Build:
RUN cd /build && go build -v


#
# Application deployment image:
#


FROM alpine:3 as run

COPY --from=build /build/rest /rest
CMD ["/rest", "-host=0.0.0.0"]
