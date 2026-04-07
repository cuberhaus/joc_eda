"""Run the C++ game engine and return the raw replay text."""
from __future__ import annotations

import subprocess
from pathlib import Path

GAME_BIN = Path("/app/game/Game")
CONFIG_FILE = Path(__file__).parent / "data" / "default.cnf"

AVAILABLE_AIS = ["Demo", "JohnnyBGood"]


def game_available() -> bool:
    return GAME_BIN.exists() and GAME_BIN.is_file()


def list_ais() -> list[str]:
    if not game_available():
        return []
    try:
        r = subprocess.run([str(GAME_BIN), "--list"], capture_output=True, text=True, timeout=5)
        return [line.strip() for line in r.stdout.strip().split("\n") if line.strip()]
    except Exception:
        return AVAILABLE_AIS


def run_game(seed: int, players: list[str] | None = None) -> str:
    if not game_available():
        raise RuntimeError("Game binary not found — C++ engine not compiled")

    if not players:
        ais = list_ais()
        num_players = 4
        players = [ais[i % len(ais)] for i in range(num_players)]

    cmd = [str(GAME_BIN), "-i", str(CONFIG_FILE), "--seed", str(seed)] + players

    result = subprocess.run(
        cmd,
        capture_output=True,
        text=True,
        timeout=30,
    )
    if result.returncode != 0:
        raise RuntimeError(f"Game exited with code {result.returncode}: {result.stderr[:500]}")
    return result.stdout
