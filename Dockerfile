FROM messense/rust-musl-cross:x86_64-musl AS builder
WORKDIR /caravan
COPY . .
RUN cargo build --release --target x86_64-unknown-linux-musl

FROM scratch
COPY --from=builder /caravan/target/x86_64-unknown-linux-musl/release/caravan /caravan
EXPOSE 8080
CMD ["/caravan"]
