FROM caddy:builder AS builder

WORKDIR /src
COPY . .

#Resolve Dependencies
RUN go mod tidy

#Compile Caddy specifying Binary Destination
RUN xcaddy build \
    --output /usr/bin/caddy \
    --with github.com/caddy-dns/vercel=/src

FROM caddy:latest

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
