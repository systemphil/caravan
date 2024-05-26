FROM rust:1.78 AS builder
WORKDIR /caravan
COPY . .
RUN cargo build --release
RUN ls -la /caravan/target/release/

FROM debian:stable-slim AS runtime
COPY --from=builder /caravan/target/release/caravan /usr/local/bin/
RUN chmod +x /usr/local/bin/caravan
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/caravan"]