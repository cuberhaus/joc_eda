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

FROM python:3.12-slim
WORKDIR /app

COPY web/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY web/backend/ ./backend/
COPY web/backend/data/ ./backend/data/
COPY default.cnf ./backend/data/default.cnf
COPY --from=frontend /build/dist ./frontend/dist/
COPY --from=engine /build/Game ./game/Game
RUN chmod +x ./game/Game

EXPOSE 8087
CMD ["uvicorn", "backend.app:app", "--host", "0.0.0.0", "--port", "8087"]
