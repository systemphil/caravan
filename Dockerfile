FROM oven/bun:latest as base
WORKDIR /usr/src/app

COPY . .
RUN bun install

ENV NODE_ENV=production
ENV PORT=3000

RUN chown -R bun:bun ./

USER bun
EXPOSE 3000/tcp
CMD [ "bun", "run", "src/index.ts" ]