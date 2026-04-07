# Joc EDA

AI player bot for the EDA (Data Structures and Algorithms) programming game competition at FIB-UPC. Players write C++ AI strategies that control warriors and builders on a grid-based board, competing against other players in real-time.

## Overview

The game takes place on a 2D grid with buildings, streets, weapons (guns, bazookas), and bonuses (money, food). Each player controls warriors and builders who move around the board, collect resources, build barricades, and fight opponents. The AI must make strategic decisions each round about movement, combat, and resource gathering.

The custom AI (`AIJohnnyMelavo.cc`) uses BFS-based pathfinding to navigate the board and prioritize objectives such as collecting weapons, hunting enemies, and gathering resources.

## Structure

```
├── AIJohnnyMelavo.cc     # Custom AI player implementation (BFS pathfinding strategy)
├── AIDemo.cc             # Demo AI player
├── Main.cc               # Game entry point
├── Board.{cc,hh}         # Board representation and logic
├── Game.{cc,hh}          # Game loop and rules
├── Player.{cc,hh}        # Player base class (to be extended by AI)
├── Structs.{cc,hh}       # Core data types (Pos, Cell, Citizen, enums)
├── Info.{cc,hh}          # Game state information
├── Action.{cc,hh}        # Action handling
├── Settings.{cc,hh}      # Game configuration
├── State.{cc,hh}         # Game state management
├── Registry.{cc,hh}      # Player registration
├── Utils.{cc,hh}         # Utility functions
├── Random.{cc,hh}        # Random number generation
├── Makefile              # Build system
├── default.cnf           # Default game configuration
├── Viewer/               # HTML-based game viewer
│   ├── viewer.html
│   ├── viewer.sh
│   └── img/              # Game sprites and assets
└── .github/workflows/    # CI pipeline
```

## Building

```bash
make
```

## Running

```bash
./Game -i default.cnf
```

Use the viewer (`Viewer/viewer.html`) to visualize game replays.

## Web App

A full-stack game viewer with live C++ engine execution, Canvas rendering, and interactive replay controls.

**Stack:** Mithril.js (Vite) + HTML5 Canvas + FastAPI backend wrapping the compiled C++ game engine

### Quick Start

```bash
# Docker (recommended)
docker compose up -d        # http://localhost:8087

# Dev mode
make web-dev                # Backend :8087, Vite dev server
```

### Features

- Live game generation — run matches with configurable AI players directly from the browser
- Replay upload and animated playback with play/pause/step/speed controls
- Canvas grid rendering with animated unit movement, buildings, weapons, and bonuses
- Keyboard shortcuts for playback control
- Player scoreboard and round-by-round stats

### Web Structure

```
web/
├── frontend/          # Mithril.js + Vite + Canvas
│   └── src/
│       ├── components/        # Board renderer, controls, player panel
│       └── styles/            # Dark theme CSS
├── backend/           # FastAPI + compiled C++ Game engine
│   └── app.py
└── requirements.txt
```

## Tech Stack

- **C++11** for game engine and AI logic
- **Mithril.js** + **Canvas** for the interactive web frontend
- **FastAPI** for the web backend (wraps the C++ engine)
- **HTML/JS** for the legacy game replay viewer (`Viewer/`)
