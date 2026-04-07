import type { GameData } from "../lib/types";

export type SourceMode = "sample" | "generate" | "upload";

export const state = {
  game: null as GameData | null,
  loading: false,
  error: "",

  currentRound: 0,
  playing: false,
  speed: 6,
  animProgress: 0,

  source: "sample" as SourceMode,
  seed: "",
  engineAvailable: false,
  availableAIs: [] as string[],
  selectedAIs: ["", "", "", ""] as string[],

  selectedPlayer: -1,
};

let animFrame = 0;
let lastTime = 0;

export function setGame(data: GameData) {
  state.game = data;
  state.currentRound = 0;
  state.playing = false;
  state.animProgress = 0;
  state.error = "";
}

export function togglePlay() {
  state.playing = !state.playing;
  if (state.playing) {
    lastTime = performance.now();
    tick();
  }
}

export function stepForward() {
  if (!state.game) return;
  if (state.currentRound < state.game.num_rounds) {
    state.currentRound++;
    state.animProgress = 0;
  }
}

export function stepBackward() {
  if (state.currentRound > 0) {
    state.currentRound--;
    state.animProgress = 0;
  }
}

export function goToStart() {
  state.currentRound = 0;
  state.animProgress = 0;
  state.playing = false;
}

export function goToEnd() {
  if (!state.game) return;
  state.currentRound = state.game.num_rounds;
  state.animProgress = 0;
  state.playing = false;
}

export function setRound(r: number) {
  state.currentRound = r;
  state.animProgress = 0;
}

function tick() {
  if (!state.playing || !state.game) return;
  const now = performance.now();
  const dt = (now - lastTime) / 1000;
  lastTime = now;

  const roundsPerSec = state.speed;
  state.animProgress += dt * roundsPerSec;

  while (state.animProgress >= 1) {
    state.animProgress -= 1;
    state.currentRound++;
    if (state.currentRound >= state.game.num_rounds) {
      state.currentRound = state.game.num_rounds;
      state.animProgress = 0;
      state.playing = false;
      const m = (window as any).__mithril;
      if (m) m.redraw();
      return;
    }
  }

  const m = (window as any).__mithril;
  if (m) m.redraw();
  animFrame = requestAnimationFrame(tick);
}

export function startPlayLoop() {
  lastTime = performance.now();
  tick();
}
