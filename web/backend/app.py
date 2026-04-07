from __future__ import annotations

import random
from pathlib import Path
from typing import Optional

from fastapi import FastAPI, HTTPException
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel

from . import parser, runner

app = FastAPI(title="ThePurge Viewer")

DATA_DIR = Path(__file__).parent / "data"
DIST_DIR = Path(__file__).parent.parent / "frontend" / "dist"

_sample_cache: dict | None = None


def _get_sample() -> dict:
    global _sample_cache
    if _sample_cache is None:
        raw = (DATA_DIR / "sample.out").read_text()
        _sample_cache = parser.parse_replay(raw)
    return _sample_cache


@app.get("/api/status")
async def status():
    return {"status": "ok", "engine_available": runner.game_available()}


class GenerateRequest(BaseModel):
    seed: Optional[int] = None
    players: Optional[list[str]] = None


@app.get("/api/sample")
async def sample():
    try:
        return _get_sample()
    except Exception as e:
        raise HTTPException(500, str(e))


@app.get("/api/ais")
async def list_ais():
    return runner.list_ais()


@app.post("/api/generate")
async def generate(req: GenerateRequest):
    if not runner.game_available():
        raise HTTPException(503, "C++ engine not available in this deployment")
    import asyncio
    seed = req.seed if req.seed is not None else random.randint(1, 999999)
    loop = asyncio.get_event_loop()
    try:
        raw = await loop.run_in_executor(
            None, lambda: runner.run_game(seed, req.players)
        )
        return parser.parse_replay(raw)
    except Exception as e:
        raise HTTPException(500, str(e))


class ParseRequest(BaseModel):
    replay_text: str


@app.post("/api/parse")
async def parse_upload(req: ParseRequest):
    try:
        return parser.parse_replay(req.replay_text)
    except Exception as e:
        raise HTTPException(400, f"Failed to parse replay: {e}")


if DIST_DIR.exists():
    app.mount("/", StaticFiles(directory=str(DIST_DIR), html=True), name="static")
