# Static docs for self-hosting (build with: docker compose build gadm-docs)
FROM peaceiris/hugo:v0.146.0-full AS hugo

WORKDIR /src
COPY . .

ARG HUGO_BASEURL=https://docs.worldlines.dev/
ENV HUGO_ENV=production

RUN hugo --minify --gc --baseURL "${HUGO_BASEURL}"

FROM nginx:1.27-alpine

COPY --from=hugo /src/public /usr/share/nginx/html

EXPOSE 80
