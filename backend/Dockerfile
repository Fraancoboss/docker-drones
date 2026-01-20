# Build (Cargo.lock v4 needs a newer Cargo)
FROM rust:1.81-bookworm AS build
WORKDIR /app
COPY . .
RUN cargo build --release

# Runtime Ubuntu 24.04 LTS (tu requisito)
FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /app/target/release/backend /app/backend
EXPOSE 8080
CMD ["/app/backend"]
