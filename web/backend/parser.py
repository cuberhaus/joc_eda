"""Parse ThePurge replay text into structured JSON-friendly dicts."""
from __future__ import annotations


def parse_replay(raw: str) -> dict:
    tokens = raw.split()
    p = 0

    def take() -> str:
        nonlocal p
        val = tokens[p]
        p += 1
        return val

    def skip(expected: str | None = None):
        nonlocal p
        if expected is not None and tokens[p] != expected:
            raise ValueError(f"Expected '{expected}', got '{tokens[p]}' at token {p}")
        p += 1

    sec_game = take() == "SecGame"
    skip("Seed")
    seed = int(take())

    game_name = take()
    if game_name != "ThePurge":
        raise ValueError("Not a ThePurge replay file")
    version = take()

    settings: dict[str, int] = {}
    setting_keys = [
        "NUM_PLAYERS", "NUM_DAYS", "NUM_ROUNDS_PER_DAY",
        "BOARD_ROWS", "BOARD_COLS",
        "NUM_INI_BUILDERS", "NUM_INI_WARRIORS",
        "NUM_INI_MONEY", "NUM_INI_FOOD", "NUM_INI_GUNS", "NUM_INI_BAZOOKAS",
        "BUILDER_INI_LIFE", "WARRIOR_INI_LIFE",
        "MONEY_POINTS", "KILL_BUILDER_POINTS", "KILL_WARRIOR_POINTS",
        "FOOD_INCR_LIFE", "LIFE_LOST_IN_ATTACK",
        "BUILDER_STRENGTH_ATTACK", "HAMMER_STRENGTH_ATTACK",
        "GUN_STRENGTH_ATTACK", "BAZOOKA_STRENGTH_ATTACK",
        "BUILDER_STRENGTH_DEMOLISH", "HAMMER_STRENGTH_DEMOLISH",
        "GUN_STRENGTH_DEMOLISH", "BAZOOKA_STRENGTH_DEMOLISH",
        "NUM_ROUNDS_REGEN_BUILDER", "NUM_ROUNDS_REGEN_WARRIOR",
        "NUM_ROUNDS_REGEN_FOOD", "NUM_ROUNDS_REGEN_MONEY",
        "NUM_ROUNDS_REGEN_WEAPON",
        "BARRICADE_RESISTANCE_STEP", "BARRICADE_MAX_RESISTANCE",
        "MAX_NUM_BARRICADES",
    ]
    for key in setting_keys:
        skip(key)
        settings[key] = int(take())

    num_players = settings["NUM_PLAYERS"]
    rows = settings["BOARD_ROWS"]
    cols = settings["BOARD_COLS"]
    num_rounds = settings["NUM_DAYS"] * settings["NUM_ROUNDS_PER_DAY"]

    skip("names")
    names = [take() for _ in range(num_players)]

    rounds = []
    for rnd in range(num_rounds + 1):
        take()  # col header line 1
        take()  # col header line 2

        grid = []
        for i in range(rows):
            label = take()
            if int(label) != i:
                raise ValueError(f"Expected row {i}, got '{label}' at token {p - 1}")
            grid.append(take())

        skip("citizens")
        n_citizens = int(take())
        skip("type"); skip("id"); skip("player")
        skip("row"); skip("column"); skip("weapon"); skip("life")

        citizens = []
        for _ in range(n_citizens):
            citizens.append({
                "type": take(),
                "id": int(take()),
                "player": int(take()),
                "row": int(take()),
                "col": int(take()),
                "weapon": take(),
                "life": int(take()),
            })

        skip("barricades")
        n_barricades = int(take())
        skip("player"); skip("row"); skip("column"); skip("resistance")

        barricades = []
        for _ in range(n_barricades):
            barricades.append({
                "player": int(take()),
                "row": int(take()),
                "col": int(take()),
                "resistance": int(take()),
            })

        skip("round")
        r_num = int(take())
        skip("day")
        is_day = int(take())

        skip("score")
        scores = [int(take()) for _ in range(num_players)]

        skip("status")
        cpu = []
        for _ in range(num_players):
            val = float(take())
            cpu.append("out" if val == -1.0 else f"{int(val * 100)}%")

        commands = []
        if rnd < num_rounds:
            skip("commands")
            n_cmds = int(take())
            for _ in range(n_cmds):
                commands.append({
                    "id": int(take()),
                    "action": take(),
                    "dir": take(),
                })

        rounds.append({
            "grid": grid,
            "citizens": citizens,
            "barricades": barricades,
            "round": r_num,
            "day": is_day,
            "scores": scores,
            "cpu": cpu,
            "commands": commands,
        })

    return {
        "sec_game": sec_game,
        "seed": seed,
        "version": version,
        "settings": settings,
        "names": names,
        "num_rounds": num_rounds,
        "rows": rows,
        "cols": cols,
        "rounds": rounds,
    }
