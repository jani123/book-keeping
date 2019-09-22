FROM alpine as ca-manager

RUN apk add --no-cache ca-certificates
RUN update-ca-certificates

FROM gobuffalo/buffalo:development as builder

ENV GO111MODULE=on

WORKDIR /build

# If dependecies already exists, do not try to fetch those, otherwise they need to be cached
ADD go.* ./
ADD vendor vendor
RUN test -f "vendor/modules.txt" || go mod download

# Do node things dependencies
ADD package.json .
RUN test -d "node_modules" || npm install

# Add sources
WORKDIR /sources
ADD . .

RUN rm -rf vendor && ln -s /build/vendor
RUN rm -rf node_modules && ln -s /build/node_modules

# Build assets
RUN npx webpack

# Build the binary
# Currently we are using self built buffalo binary, but remove it after this
# pull request has gone through https://github.com/gobuffalo/buffalo/pull/1675
RUN buffalo build --static --skip-build-deps --mod=vendor -o /bin/app

FROM scratch

#RUN apk add --no-cache curl

WORKDIR /bin/
COPY --from=builder /bin/app .
COPY --from=builder /go/bin/buffalo /bin/buffalo
COPY --from=builder /sources/public /public
COPY --from=ca-manager /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0
ENV PORT=3000

EXPOSE 3000

ENTRYPOINT ["/bin/app"]
