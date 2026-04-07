FROM node:22-slim AS frontend
WORKDIR /build
COPY web/frontend/package.json web/frontend/package-lock.json ./
RUN npm ci
COPY web/frontend/ ./
RUN npm run build

FROM gcc:13-bookworm AS engine
WORKDIR /build
COPY *.cc *.hh Makefile ./
RUN touch Makefile.deps && make DUMMY_OBJ= EXTRA_OBJ= Game

FROM golang:1.22-bookworm AS go-build
WORKDIR /src
COPY web/backend-go/go.mod ./
RUN go mod download
COPY web/backend-go/*.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=engine /usr/local/lib64/libstdc++.so.6* /usr/local/lib64/
ENV LD_LIBRARY_PATH=/usr/local/lib64
COPY --from=engine /build/Game ./game/Game
RUN chmod +x ./game/Game
COPY --from=go-build /src/server ./server
COPY --from=frontend /build/dist ./frontend/dist/
COPY web/backend/data/ ./data/
COPY default.cnf ./data/default.cnf

ENV PORT=8087
ENV GAME_BIN=/app/game/Game
ENV CONFIG_FILE=/app/data/default.cnf
ENV DATA_DIR=/app/data
ENV STATIC_DIR=/app/frontend/dist
EXPOSE 8087

CMD ["./server"]
