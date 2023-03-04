FROM denoland/deno:1.31.1

WORKDIR /app
COPY index.ts .
RUN deno cache index.ts

ENTRYPOINT [ "deno", "run", "--allow-net", "index.ts" ]
