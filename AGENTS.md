# joc_eda

FIB-UPC EDA (Data Structures and Algorithms) programming-game project. Contains a graded C++ AI player, an HTML5 Canvas + Mithril.js replay viewer, and a Go backend that wraps the compiled C++ engine to serve live matches.

## Architecture

- **C++ game engine + AI** (repo root): course-provided framework (`Game`, `Board`, `Player`, etc.) plus the custom AI `AIJohnnyMelavo.cc` (BFS-based pathfinding). C++11, built via `Makefile`. This is the graded deliverable.
- **Web frontend** ([web/frontend/](web/frontend/)): Mithril.js + Vite + HTML5 Canvas. Renders replays, exposes play/pause/step/speed controls and a match-launch form.
- **Web backend** ([web/backend-go/](web/backend-go/)): Go stdlib `net/http` server. Spawns the compiled `Game` binary as a subprocess, parses its token-based replay output, and serves the SPA. Default port `8087`.
- **Legacy viewer** ([Viewer/](Viewer/)): standalone HTML replay viewer kept for parity with the course tooling.

## Build and Test

```bash
# C++ engine + AI (produces ./Game)
make
./Game -i default.cnf -o match.res     # run a match with the default config

# Web app (Docker, recommended)
docker compose up -d                   # http://localhost:8087

# Web app (dev)
make web-install                       # one-time: npm install in web/frontend
make web-dev                           # Go backend on :8087 + Vite dev server
make web-build                         # production frontend bundle
```

The Go backend expects `GAME_BIN`, `CONFIG_FILE`, `DATA_DIR`, `STATIC_DIR` env vars (see `web-dev` target in [Makefile](Makefile)).

## Pitfalls

- `AIJohnnyMelavo.cc` and the surrounding C++ framework are the graded course artifact — treat as **frozen**. Do not refactor framework files (`Game.*`, `Board.*`, `Player.*`, `Structs.*`, etc.); only the `AI*.cc` strategy file is fair game.
- The Makefile auto-detects all `AI*.cc` files. Adding a new AI requires no Makefile change but the file must follow the `Player` interface and call the `Registry` macro.
- Web frontend/backend are auxiliary tooling and can evolve freely, but the backend depends on the exact stdout token format produced by `./Game` — keep `web/backend-go/parser.go` in sync if the engine output ever changes.
- `default.cnf` is the canonical match config; `default-fixed.cnf` is used for deterministic regression runs.

See [README.md](README.md).
