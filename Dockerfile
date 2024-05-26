FROM rust:1.78 AS builder
WORKDIR /caravan
COPY . .
RUN cargo build --release

FROM debian:stable-slim AS runtime
COPY --from=builder /caravan/target/release/caravan /usr/local/bin/
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/caravan"]
