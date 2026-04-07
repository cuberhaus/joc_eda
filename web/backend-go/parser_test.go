package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestTokenReaderTake(t *testing.T) {
	tr := &tokenReader{tokens: []string{"a", "b", "c"}}

	v, err := tr.take()
	if err != nil || v != "a" {
		t.Fatalf("expected 'a', got %q err=%v", v, err)
	}
	v, _ = tr.take()
	if v != "b" {
		t.Fatalf("expected 'b', got %q", v)
	}
}

func TestTokenReaderTakeInt(t *testing.T) {
	tr := &tokenReader{tokens: []string{"42", "-3", "abc"}}

	v, err := tr.takeInt()
	if err != nil || v != 42 {
		t.Fatalf("expected 42, got %d err=%v", v, err)
	}
	v, err = tr.takeInt()
	if err != nil || v != -3 {
		t.Fatalf("expected -3, got %d err=%v", v, err)
	}
	_, err = tr.takeInt()
	if err == nil {
		t.Fatal("expected error for non-integer")
	}
}

func TestTokenReaderTakeFloat(t *testing.T) {
	tr := &tokenReader{tokens: []string{"3.14", "-1.0"}}

	v, err := tr.takeFloat()
	if err != nil || v != 3.14 {
		t.Fatalf("expected 3.14, got %f err=%v", v, err)
	}
	v, err = tr.takeFloat()
	if err != nil || v != -1.0 {
		t.Fatalf("expected -1.0, got %f err=%v", v, err)
	}
}

func TestTokenReaderExhausted(t *testing.T) {
	tr := &tokenReader{tokens: []string{"x"}}
	tr.take()
	_, err := tr.take()
	if err == nil {
		t.Fatal("expected error when exhausted")
	}
}

func TestTokenReaderSkip(t *testing.T) {
	tr := &tokenReader{tokens: []string{"Seed", "42"}}
	if err := tr.skip("Seed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := tr.skip("wrong"); err == nil {
		t.Fatal("expected error for wrong token")
	}
}

func TestTokenReaderSkipEmpty(t *testing.T) {
	tr := &tokenReader{tokens: []string{"anything"}}
	if err := tr.skip(""); err != nil {
		t.Fatal("skip with empty expected should accept anything")
	}
}

func buildTestReplay(numPlayers, rows, cols, days, roundsPerDay int) string {
	var parts []string
	parts = append(parts, "Game", "Seed", "12345", "ThePurge", "1.0")

	settings := map[string]int{
		"NUM_PLAYERS": numPlayers, "NUM_DAYS": days, "NUM_ROUNDS_PER_DAY": roundsPerDay,
		"BOARD_ROWS": rows, "BOARD_COLS": cols,
		"NUM_INI_BUILDERS": 2, "NUM_INI_WARRIORS": 2,
		"NUM_INI_MONEY": 10, "NUM_INI_FOOD": 10, "NUM_INI_GUNS": 1, "NUM_INI_BAZOOKAS": 0,
		"BUILDER_INI_LIFE": 50, "WARRIOR_INI_LIFE": 80,
		"MONEY_POINTS": 5, "KILL_BUILDER_POINTS": 10, "KILL_WARRIOR_POINTS": 15,
		"FOOD_INCR_LIFE": 20, "LIFE_LOST_IN_ATTACK": 10,
		"BUILDER_STRENGTH_ATTACK": 3, "HAMMER_STRENGTH_ATTACK": 5,
		"GUN_STRENGTH_ATTACK": 8, "BAZOOKA_STRENGTH_ATTACK": 15,
		"BUILDER_STRENGTH_DEMOLISH": 3, "HAMMER_STRENGTH_DEMOLISH": 5,
		"GUN_STRENGTH_DEMOLISH": 8, "BAZOOKA_STRENGTH_DEMOLISH": 15,
		"NUM_ROUNDS_REGEN_BUILDER": 5, "NUM_ROUNDS_REGEN_WARRIOR": 5,
		"NUM_ROUNDS_REGEN_FOOD": 3, "NUM_ROUNDS_REGEN_MONEY": 3,
		"NUM_ROUNDS_REGEN_WEAPON": 4,
		"BARRICADE_RESISTANCE_STEP": 2, "BARRICADE_MAX_RESISTANCE": 10,
		"MAX_NUM_BARRICADES": 3,
	}
	for _, key := range settingKeys {
		parts = append(parts, key, fmt.Sprintf("%d", settings[key]))
	}

	parts = append(parts, "names")
	for i := 0; i < numPlayers; i++ {
		parts = append(parts, fmt.Sprintf("P%d", i))
	}

	totalRounds := days * roundsPerDay
	gridRow := strings.Repeat(".", cols)
	for rnd := 0; rnd <= totalRounds; rnd++ {
		parts = append(parts, "hdr1", "hdr2")
		for r := 0; r < rows; r++ {
			parts = append(parts, fmt.Sprintf("r%d", r), gridRow)
		}
		parts = append(parts, "citizens", "0",
			"type", "id", "player", "row", "column", "weapon", "life")
		parts = append(parts, "barricades", "0",
			"player", "row", "column", "resistance")
		parts = append(parts, "round", fmt.Sprintf("%d", rnd), "day", "0")
		parts = append(parts, "score")
		for i := 0; i < numPlayers; i++ {
			parts = append(parts, "0")
		}
		parts = append(parts, "status")
		for i := 0; i < numPlayers; i++ {
			parts = append(parts, "0.5")
		}
		if rnd < totalRounds {
			parts = append(parts, "commands", "0")
		}
	}

	return strings.Join(parts, " ")
}

func buildReplayWithCitizens() string {
	var parts []string
	parts = append(parts, "Game", "Seed", "12345", "ThePurge", "1.0")

	settings := map[string]int{
		"NUM_PLAYERS": 2, "NUM_DAYS": 1, "NUM_ROUNDS_PER_DAY": 1,
		"BOARD_ROWS": 3, "BOARD_COLS": 3,
		"NUM_INI_BUILDERS": 2, "NUM_INI_WARRIORS": 2,
		"NUM_INI_MONEY": 10, "NUM_INI_FOOD": 10, "NUM_INI_GUNS": 1, "NUM_INI_BAZOOKAS": 0,
		"BUILDER_INI_LIFE": 50, "WARRIOR_INI_LIFE": 80,
		"MONEY_POINTS": 5, "KILL_BUILDER_POINTS": 10, "KILL_WARRIOR_POINTS": 15,
		"FOOD_INCR_LIFE": 20, "LIFE_LOST_IN_ATTACK": 10,
		"BUILDER_STRENGTH_ATTACK": 3, "HAMMER_STRENGTH_ATTACK": 5,
		"GUN_STRENGTH_ATTACK": 8, "BAZOOKA_STRENGTH_ATTACK": 15,
		"BUILDER_STRENGTH_DEMOLISH": 3, "HAMMER_STRENGTH_DEMOLISH": 5,
		"GUN_STRENGTH_DEMOLISH": 8, "BAZOOKA_STRENGTH_DEMOLISH": 15,
		"NUM_ROUNDS_REGEN_BUILDER": 5, "NUM_ROUNDS_REGEN_WARRIOR": 5,
		"NUM_ROUNDS_REGEN_FOOD": 3, "NUM_ROUNDS_REGEN_MONEY": 3,
		"NUM_ROUNDS_REGEN_WEAPON": 4,
		"BARRICADE_RESISTANCE_STEP": 2, "BARRICADE_MAX_RESISTANCE": 10,
		"MAX_NUM_BARRICADES": 3,
	}
	for _, key := range settingKeys {
		parts = append(parts, key, fmt.Sprintf("%d", settings[key]))
	}

	parts = append(parts, "names", "P0", "P1")

	for rnd := 0; rnd <= 1; rnd++ {
		parts = append(parts, "hdr1", "hdr2")
		for r := 0; r < 3; r++ {
			parts = append(parts, fmt.Sprintf("r%d", r), "...")
		}

		if rnd == 0 {
			parts = append(parts, "citizens", "1",
				"type", "id", "player", "row", "column", "weapon", "life",
				"warrior", "7", "0", "1", "2", "gun", "80")
		} else {
			parts = append(parts, "citizens", "0",
				"type", "id", "player", "row", "column", "weapon", "life")
		}

		parts = append(parts, "barricades", "0",
			"player", "row", "column", "resistance")
		parts = append(parts, "round", fmt.Sprintf("%d", rnd), "day", "0")
		parts = append(parts, "score", "0", "0")
		parts = append(parts, "status", "0.5", "0.5")
		if rnd < 1 {
			parts = append(parts, "commands", "0")
		}
	}

	return strings.Join(parts, " ")
}

func TestParseReplayMinimal(t *testing.T) {
	raw := buildTestReplay(2, 3, 3, 1, 1)
	replay, err := parseReplay(raw)
	if err != nil {
		t.Fatalf("parseReplay failed: %v", err)
	}

	if replay.Seed != 12345 {
		t.Errorf("expected seed 12345, got %d", replay.Seed)
	}
	if replay.SecGame {
		t.Error("expected SecGame=false")
	}
	if replay.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", replay.Version)
	}
	if len(replay.Names) != 2 {
		t.Errorf("expected 2 players, got %d", len(replay.Names))
	}
	if replay.Rows != 3 || replay.Cols != 3 {
		t.Errorf("expected 3x3 board, got %dx%d", replay.Rows, replay.Cols)
	}
	if replay.NumRounds != 1 {
		t.Errorf("expected 1 round, got %d", replay.NumRounds)
	}
	if len(replay.Rounds) != 2 {
		t.Errorf("expected 2 round entries (0..numRounds), got %d", len(replay.Rounds))
	}
}

func TestParseReplaySettings(t *testing.T) {
	raw := buildTestReplay(4, 5, 5, 2, 3)
	replay, err := parseReplay(raw)
	if err != nil {
		t.Fatalf("parseReplay failed: %v", err)
	}
	if replay.Settings["NUM_PLAYERS"] != 4 {
		t.Errorf("expected NUM_PLAYERS=4, got %d", replay.Settings["NUM_PLAYERS"])
	}
	if replay.Settings["BOARD_ROWS"] != 5 {
		t.Errorf("expected BOARD_ROWS=5, got %d", replay.Settings["BOARD_ROWS"])
	}
	if replay.NumRounds != 6 {
		t.Errorf("expected 6 rounds (2*3), got %d", replay.NumRounds)
	}
}

func TestParseReplaySecGame(t *testing.T) {
	raw := buildTestReplay(2, 3, 3, 1, 1)
	raw = "SecGame " + raw[len("Game "):]
	replay, err := parseReplay(raw)
	if err != nil {
		t.Fatalf("parseReplay failed: %v", err)
	}
	if !replay.SecGame {
		t.Error("expected SecGame=true")
	}
}

func TestParseReplayWithCitizens(t *testing.T) {
	raw := buildReplayWithCitizens()
	replay, err := parseReplay(raw)
	if err != nil {
		t.Fatalf("parseReplay failed: %v", err)
	}
	if len(replay.Rounds[0].Citizens) != 1 {
		t.Errorf("expected 1 citizen in round 0, got %d", len(replay.Rounds[0].Citizens))
	}
	c := replay.Rounds[0].Citizens[0]
	if c.Type != "warrior" {
		t.Errorf("expected warrior, got %s", c.Type)
	}
	if c.ID != 7 {
		t.Errorf("expected id 7, got %d", c.ID)
	}
	if c.Player != 0 {
		t.Errorf("expected player 0, got %d", c.Player)
	}
	if c.Life != 80 {
		t.Errorf("expected life 80, got %d", c.Life)
	}
	if c.Weapon != "gun" {
		t.Errorf("expected weapon gun, got %s", c.Weapon)
	}
	if len(replay.Rounds[1].Citizens) != 0 {
		t.Errorf("expected 0 citizens in round 1, got %d", len(replay.Rounds[1].Citizens))
	}
}

func TestParseReplayEmpty(t *testing.T) {
	_, err := parseReplay("")
	if err == nil {
		t.Error("expected error on empty input")
	}
}

func TestParseReplayInvalid(t *testing.T) {
	_, err := parseReplay("not a valid replay at all")
	if err == nil {
		t.Error("expected error on invalid input")
	}
}

func TestCPUStatusFormatting(t *testing.T) {
	raw := buildTestReplay(2, 3, 3, 1, 1)
	replay, err := parseReplay(raw)
	if err != nil {
		t.Fatalf("parseReplay failed: %v", err)
	}
	for _, rnd := range replay.Rounds {
		for _, cpu := range rnd.CPU {
			if cpu != "50%" && cpu != "out" {
				t.Errorf("unexpected cpu format: %q", cpu)
			}
		}
	}
}

func TestCPUStatusOut(t *testing.T) {
	raw := buildTestReplay(2, 3, 3, 1, 1)
	raw = strings.Replace(raw, "status 0.5 0.5", "status -1.0 0.5", 1)
	replay, err := parseReplay(raw)
	if err != nil {
		t.Fatalf("parseReplay failed: %v", err)
	}
	if replay.Rounds[0].CPU[0] != "out" {
		t.Errorf("expected 'out' for -1.0, got %q", replay.Rounds[0].CPU[0])
	}
	if replay.Rounds[0].CPU[1] != "50%" {
		t.Errorf("expected '50%%', got %q", replay.Rounds[0].CPU[1])
	}
}

func TestParseReplayGridDimensions(t *testing.T) {
	raw := buildTestReplay(2, 4, 5, 1, 1)
	replay, err := parseReplay(raw)
	if err != nil {
		t.Fatalf("parseReplay failed: %v", err)
	}
	for _, rnd := range replay.Rounds {
		if len(rnd.Grid) != 4 {
			t.Errorf("expected 4 grid rows, got %d", len(rnd.Grid))
		}
		for _, row := range rnd.Grid {
			if len(row) != 5 {
				t.Errorf("expected grid row length 5, got %d", len(row))
			}
		}
	}
}
