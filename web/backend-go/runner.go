package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	gameBin    = envOr("GAME_BIN", "/app/game/Game")
	configFile = envOr("CONFIG_FILE", "data/default.cnf")

	fallbackAIs = []string{"Demo", "JohnnyBGood"}
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func gameAvailable() bool {
	info, err := os.Stat(gameBin)
	return err == nil && !info.IsDir()
}

func listAIs() []string {
	if !gameAvailable() {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, gameBin, "--list").Output()
	if err != nil {
		return fallbackAIs
	}

	var ais []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if s := strings.TrimSpace(line); s != "" {
			ais = append(ais, s)
		}
	}
	if len(ais) == 0 {
		return fallbackAIs
	}
	return ais
}

func runGame(seed int, players []string) (string, error) {
	if !gameAvailable() {
		return "", fmt.Errorf("game binary not found — C++ engine not compiled")
	}

	if len(players) == 0 {
		ais := listAIs()
		numPlayers := 4
		players = make([]string, numPlayers)
		for i := range players {
			players[i] = ais[i%len(ais)]
		}
	}

	args := []string{"-i", configFile, "--seed", strconv.Itoa(seed)}
	args = append(args, players...)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, gameBin, args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)
			if len(stderr) > 500 {
				stderr = stderr[:500]
			}
			return "", fmt.Errorf("game exited with code %d: %s", exitErr.ExitCode(), stderr)
		}
		return "", err
	}
	return string(out), nil
}
