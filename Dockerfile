FROM rust:1.78 AS builder
WORKDIR /caravan
COPY . .
RUN cargo build --release

FROM debian:stable-slim AS runtime
COPY --from=builder /caravan/target/release/caravan /caravan
EXPOSE 8080
CMD ["/usr/local/bin/caravan"]
