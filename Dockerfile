FROM golang:latest
MAINTAINER Chris Aumann

# Install glide Go dependency manager
RUN curl --silent -L https://github.com/Masterminds/glide/releases/download/v0.13.1/glide-v0.13.1-linux-amd64.tar.gz |tar -xz --strip-components=1 -C /usr/bin

# Create application directory and switch to it
RUN mkdir -p /go/src/github.com/cndy-store/analytics
WORKDIR /go/src/github.com/cndy-store/analytics

# Get dependencies. Copy glide.* before the actual code, so dependencies are only refeched on changes.
# Otherwise, dependencies are re-fetched everytime we're changing the code.
ADD glide.yaml /go/src/github.com/cndy-store/analytics/glide.yaml
ADD glide.yaml /go/src/github.com/cndy-store/analytics/glide.lock
RUN glide install

ADD . /go/src/github.com/cndy-store/analytics
RUN go build

# Run in production mode
# ENV GIN_MODE release

CMD ["sh", "-c", "./wait-for-it.sh ${PGHOST}:5432 -- ./analytics"]
