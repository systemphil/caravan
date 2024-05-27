FROM rust:1.78 AS builder
WORKDIR /build
COPY . .
# RUN cargo install --path .
RUN cargo build --release

# FROM debian:bullseye-slim
# RUN apt-get update && apt-get install -y libnss3 libgconf-2-4 libatk1.0-0 libatk-bridge2.0-0 libgdk-pixbuf2.0-0 libgtk-3-0 libgbm-dev libnss3-dev libgdk-pixbuf2.0-dev libgtk-3-dev libxss-dev && rm -rf /var/lib/apt/lists/*
# COPY --from=builder /build/target/release/caravan /usr/local/bin/
# # RUN apt-get update
# # RUN apt-get install -y libnss3
# # RUN chmod +x /usr/local/bin/caravan

EXPOSE 8080
CMD ["/build/target/release/caravan"]