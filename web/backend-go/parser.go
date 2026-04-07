package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Citizen struct {
	Type   string `json:"type"`
	ID     int    `json:"id"`
	Player int    `json:"player"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Weapon string `json:"weapon"`
	Life   int    `json:"life"`
}

type Barricade struct {
	Player     int `json:"player"`
	Row        int `json:"row"`
	Col        int `json:"col"`
	Resistance int `json:"resistance"`
}

type Command struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
	Dir    string `json:"dir"`
}

type Round struct {
	Grid       []string    `json:"grid"`
	Citizens   []Citizen   `json:"citizens"`
	Barricades []Barricade `json:"barricades"`
	Round      int         `json:"round"`
	Day        int         `json:"day"`
	Scores     []int       `json:"scores"`
	CPU        []string    `json:"cpu"`
	Commands   []Command   `json:"commands"`
}

type Replay struct {
	SecGame   bool           `json:"sec_game"`
	Seed      int            `json:"seed"`
	Version   string         `json:"version"`
	Settings  map[string]int `json:"settings"`
	Names     []string       `json:"names"`
	NumRounds int            `json:"num_rounds"`
	Rows      int            `json:"rows"`
	Cols      int            `json:"cols"`
	Rounds    []Round        `json:"rounds"`
}

type tokenReader struct {
	tokens []string
	pos    int
}

func (t *tokenReader) take() (string, error) {
	if t.pos >= len(t.tokens) {
		return "", fmt.Errorf("unexpected end of tokens at position %d", t.pos)
	}
	val := t.tokens[t.pos]
	t.pos++
	return val, nil
}

func (t *tokenReader) takeInt() (int, error) {
	s, err := t.take()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func (t *tokenReader) takeFloat() (float64, error) {
	s, err := t.take()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

func (t *tokenReader) skip(expected string) error {
	s, err := t.take()
	if err != nil {
		return err
	}
	if expected != "" && s != expected {
		return fmt.Errorf("expected %q, got %q at token %d", expected, s, t.pos-1)
	}
	return nil
}

var settingKeys = []string{
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
}

func parseReplay(raw string) (*Replay, error) {
	tr := &tokenReader{tokens: strings.Fields(raw)}

	first, err := tr.take()
	if err != nil {
		return nil, err
	}
	secGame := first == "SecGame"

	if err := tr.skip("Seed"); err != nil {
		return nil, err
	}
	seed, err := tr.takeInt()
	if err != nil {
		return nil, err
	}

	gameName, err := tr.take()
	if err != nil {
		return nil, err
	}
	if gameName != "ThePurge" {
		return nil, fmt.Errorf("not a ThePurge replay file")
	}
	version, err := tr.take()
	if err != nil {
		return nil, err
	}

	settings := make(map[string]int, len(settingKeys))
	for _, key := range settingKeys {
		if err := tr.skip(key); err != nil {
			return nil, err
		}
		v, err := tr.takeInt()
		if err != nil {
			return nil, err
		}
		settings[key] = v
	}

	numPlayers := settings["NUM_PLAYERS"]
	rows := settings["BOARD_ROWS"]
	cols := settings["BOARD_COLS"]
	numRounds := settings["NUM_DAYS"] * settings["NUM_ROUNDS_PER_DAY"]

	if err := tr.skip("names"); err != nil {
		return nil, err
	}
	names := make([]string, numPlayers)
	for i := range names {
		names[i], err = tr.take()
		if err != nil {
			return nil, err
		}
	}

	rounds := make([]Round, 0, numRounds+1)
	for rnd := 0; rnd <= numRounds; rnd++ {
		tr.take() // col header 1
		tr.take() // col header 2

		grid := make([]string, rows)
		for i := 0; i < rows; i++ {
			tr.take() // row label
			grid[i], err = tr.take()
			if err != nil {
				return nil, err
			}
		}

		if err := tr.skip("citizens"); err != nil {
			return nil, err
		}
		nCitizens, err := tr.takeInt()
		if err != nil {
			return nil, err
		}
		for _, h := range []string{"type", "id", "player", "row", "column", "weapon", "life"} {
			tr.skip(h)
		}

		citizens := make([]Citizen, nCitizens)
		for i := 0; i < nCitizens; i++ {
			cType, _ := tr.take()
			cID, _ := tr.takeInt()
			cPlayer, _ := tr.takeInt()
			cRow, _ := tr.takeInt()
			cCol, _ := tr.takeInt()
			cWeapon, _ := tr.take()
			cLife, _ := tr.takeInt()
			citizens[i] = Citizen{cType, cID, cPlayer, cRow, cCol, cWeapon, cLife}
		}

		if err := tr.skip("barricades"); err != nil {
			return nil, err
		}
		nBarricades, err := tr.takeInt()
		if err != nil {
			return nil, err
		}
		for _, h := range []string{"player", "row", "column", "resistance"} {
			tr.skip(h)
		}

		barricades := make([]Barricade, nBarricades)
		for i := 0; i < nBarricades; i++ {
			bPlayer, _ := tr.takeInt()
			bRow, _ := tr.takeInt()
			bCol, _ := tr.takeInt()
			bRes, _ := tr.takeInt()
			barricades[i] = Barricade{bPlayer, bRow, bCol, bRes}
		}

		tr.skip("round")
		roundNum, _ := tr.takeInt()
		tr.skip("day")
		isDay, _ := tr.takeInt()

		tr.skip("score")
		scores := make([]int, numPlayers)
		for i := range scores {
			scores[i], _ = tr.takeInt()
		}

		tr.skip("status")
		cpu := make([]string, numPlayers)
		for i := range cpu {
			val, _ := tr.takeFloat()
			if val == -1.0 {
				cpu[i] = "out"
			} else {
				cpu[i] = fmt.Sprintf("%d%%", int(val*100))
			}
		}

		var commands []Command
		if rnd < numRounds {
			tr.skip("commands")
			nCmds, _ := tr.takeInt()
			commands = make([]Command, nCmds)
			for i := 0; i < nCmds; i++ {
				cID, _ := tr.takeInt()
				cAction, _ := tr.take()
				cDir, _ := tr.take()
				commands[i] = Command{cID, cAction, cDir}
			}
		} else {
			commands = []Command{}
		}

		rounds = append(rounds, Round{
			Grid:       grid,
			Citizens:   citizens,
			Barricades: barricades,
			Round:      roundNum,
			Day:        isDay,
			Scores:     scores,
			CPU:        cpu,
			Commands:   commands,
		})
	}

	return &Replay{
		SecGame:   secGame,
		Seed:      seed,
		Version:   version,
		Settings:  settings,
		Names:     names,
		NumRounds: numRounds,
		Rows:      rows,
		Cols:      cols,
		Rounds:    rounds,
	}, nil
}
