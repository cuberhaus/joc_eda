import m from "mithril";
import { GameCanvas } from "./canvas";
import {
  state,
  togglePlay,
  stepForward,
  stepBackward,
  goToStart,
  goToEnd,
  setRound,
  setGame,
  startPlayLoop,
} from "./state";
import { getSample, generateGame, parseReplay, getStatus, getAIs } from "../lib/api";
import type { Citizen } from "../lib/types";

const PLAYER_COLORS = ["var(--player0)", "var(--player1)", "var(--player2)", "var(--player3)"];
const WEAPON_NAMES: Record<string, string> = { n: "None", h: "Hammer", g: "Gun", b: "Bazooka" };
const TYPE_NAMES: Record<string, string> = { b: "Builder", w: "Warrior" };

let canvasInst: GameCanvas | null = null;
let fileInput: HTMLInputElement | null = null;

function renderCanvas(vnode: m.VnodeDOM) {
  if (!state.game) return;
  const canvas = vnode.dom.querySelector("canvas") as HTMLCanvasElement;
  if (!canvas) return;
  if (!canvasInst) canvasInst = new GameCanvas(canvas);
  canvasInst.render(state.game, state.currentRound, state.animProgress);
}

async function loadSample() {
  state.loading = true;
  state.error = "";
  m.redraw();
  try {
    const data = await getSample();
    setGame(data);
  } catch (e: any) {
    state.error = e.message;
  }
  state.loading = false;
  m.redraw();
}

async function generate() {
  state.loading = true;
  state.error = "";
  m.redraw();
  try {
    const seed = state.seed ? parseInt(state.seed) : undefined;
    const players = state.selectedAIs.filter(Boolean).length === 4
      ? state.selectedAIs
      : undefined;
    const data = await generateGame(seed, players);
    setGame(data);
  } catch (e: any) {
    state.error = e.message;
  }
  state.loading = false;
  m.redraw();
}

function handleFile(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (!file) return;
  state.loading = true;
  state.error = "";
  m.redraw();
  const reader = new FileReader();
  reader.onload = async () => {
    try {
      const data = await parseReplay(reader.result as string);
      setGame(data);
    } catch (err: any) {
      state.error = err.message;
    }
    state.loading = false;
    m.redraw();
  };
  reader.readAsText(file);
}

function onPlay() {
  togglePlay();
  if (state.playing) startPlayLoop();
}

function LeftPanel(): m.Component {
  return {
    oninit() {
      getStatus().then((s) => {
        state.engineAvailable = s.engine_available;
        if (s.engine_available) {
          getAIs().then((ais) => {
            state.availableAIs = ais;
            state.selectedAIs = ais.length >= 4
              ? [ais[0], ais[0], ais[1 % ais.length], ais[1 % ais.length]]
              : ais.concat(ais).slice(0, 4);
            m.redraw();
          });
        }
        m.redraw();
      });
    },
    view() {
      return m("aside.panel.panel-left", [
        m("div.section-title", "Source"),
        m("div.source-btns", [
          m("button.btn", { class: state.source === "sample" ? "active" : "", onclick: () => { state.source = "sample"; } }, "Sample"),
          m("button.btn", { class: state.source === "generate" ? "active" : "", onclick: () => { state.source = "generate"; } }, "Generate"),
          m("button.btn", { class: state.source === "upload" ? "active" : "", onclick: () => { state.source = "upload"; } }, "Upload"),
        ]),

        state.source === "sample" && m("div", [
          m("button.btn.btn-primary", {
            style: "width:100%",
            onclick: loadSample,
            disabled: state.loading,
          }, state.loading ? "Loading..." : "Load Sample Game"),
        ]),

        state.source === "generate" && m("div", [
          m("div.form-group", [
            m("label.form-label", "Seed (optional)"),
            m("input.form-input[type=number]", {
              placeholder: "Random",
              value: state.seed,
              oninput: (e: InputEvent) => { state.seed = (e.target as HTMLInputElement).value; },
            }),
          ]),
          state.availableAIs.length > 0 && m("div.form-group", [
            m("label.form-label", "Players"),
            ...[0, 1, 2, 3].map(i =>
              m("select.form-select", {
                style: `margin-bottom:0.15rem;border-left:3px solid ${["var(--player0)","var(--player1)","var(--player2)","var(--player3)"][i]}`,
                value: state.selectedAIs[i],
                onchange: (e: Event) => { state.selectedAIs[i] = (e.target as HTMLSelectElement).value; },
              },
                state.availableAIs.map(ai => m("option", { value: ai }, ai))
              )
            ),
          ]),
          m("button.btn.btn-primary", {
            style: "width:100%",
            onclick: generate,
            disabled: state.loading || !state.engineAvailable,
          }, state.loading ? "Generating..." : !state.engineAvailable ? "Engine N/A" : "Generate Game"),
          !state.engineAvailable && m("div", { style: "font-size:0.65rem;color:var(--text-muted);margin-top:0.3rem" },
            "C++ engine not compiled in this deployment. Use Docker build."),
        ]),

        state.source === "upload" && m("div", [
          m("button.btn", {
            style: "width:100%",
            onclick: () => fileInput?.click(),
          }, "Choose .out File"),
          m("input[type=file][accept=.out,.txt]", {
            style: "display:none",
            oncreate: (v: m.VnodeDOM) => { fileInput = v.dom as HTMLInputElement; },
            onchange: handleFile,
          }),
        ]),

        state.error && m("div", { style: "color:var(--danger);font-size:0.72rem;margin-top:0.4rem;word-break:break-word" }, state.error),

        state.game && m("div", [
          m("div.section-title", { style: "margin-top:0.8rem" }, "Game Info"),
          m("div", { style: "font-size:0.72rem;color:var(--text-secondary)" }, [
            m("div", `Seed: ${state.game.seed}`),
            m("div", `Grid: ${state.game.rows} × ${state.game.cols}`),
            m("div", `Rounds: ${state.game.num_rounds}`),
            m("div", `Players: ${state.game.names.join(", ")}`),
          ]),

          m("div.section-title", { style: "margin-top:0.8rem" }, "Legend"),
          m("div", { style: "font-size:0.68rem;color:var(--text-muted);display:grid;grid-template-columns:auto 1fr;gap:0.15rem 0.5rem;align-items:center" }, [
            m("div", { style: "width:10px;height:10px;border-radius:50%;background:#22c55e" }), m("span", "Warrior"),
            m("div", { style: "width:10px;height:10px;background:#22c55e" }), m("span", "Builder"),
            m("div", { style: "width:10px;height:10px;border-radius:50%;background:#f59e0b" }), m("span", "Food"),
            m("div", { style: "width:10px;height:10px;background:#fbbf24" }), m("span", "Money"),
            m("div", { style: "width:0;height:0;border-left:5px solid transparent;border-right:5px solid transparent;border-bottom:10px solid #94a3b8" }), m("span", "Gun"),
            m("div", { style: "width:10px;height:10px;background:#f97316;clip-path:polygon(50% 0%,100% 38%,82% 100%,18% 100%,0% 38%)" }), m("span", "Bazooka"),
            m("div", { style: "width:10px;height:10px;background:#2a2a3e;border:1px solid #3a3a50" }), m("span", "Building"),
            m("div", { style: "width:10px;height:10px;border:1px dashed #888" }), m("span", "Barricade"),
          ]),
        ]),
      ]);
    },
  };
}

function RightPanel(): m.Component {
  return {
    view() {
      if (!state.game) {
        return m("aside.panel.panel-right", m("div", { style: "color:var(--text-muted);font-size:0.75rem;padding:1rem" }, "Load a game to see scores and units."));
      }

      const rd = state.game.rounds[state.currentRound];
      const maxScore = Math.max(1, ...rd.scores);

      return m("aside.panel.panel-right", [
        m("div.section-title", "Scores"),
        state.game.names.map((name, i) =>
          m("div.player-card", {
            style: state.selectedPlayer === i ? `border-color:${PLAYER_COLORS[i].replace("var(--player" + i + ")", "")}` : "",
            onclick: () => { state.selectedPlayer = state.selectedPlayer === i ? -1 : i; },
          }, [
            m("div.player-dot", { style: `background:${PLAYER_COLORS[i]}` }),
            m("span.player-name", name),
            m("span.player-score", rd.scores[i]),
            m("span.player-cpu", rd.cpu[i]),
          ])
        ),

        m("div.score-bar-wrap",
          state.game.names.map((_, i) =>
            m("div.score-bar-row", [
              m("div.bar", { style: `width:${(rd.scores[i] / maxScore) * 100}%;background:${PLAYER_COLORS[i]}` }),
              m("span.val", rd.scores[i]),
            ])
          )
        ),

        m("div.section-title", { style: "margin-top:0.6rem" }, `Units (${rd.citizens.length})`),
        m("div.unit-list",
          rd.citizens
            .filter((c: Citizen) => state.selectedPlayer === -1 || c.player === state.selectedPlayer)
            .map((c: Citizen) =>
              m("div.unit-row", [
                m("span.unit-type", { style: `color:${PLAYER_COLORS[c.player]}` }, c.type === "b" ? "B" : "W"),
                m("span.unit-id", `#${c.id}`),
                m("span", { style: "flex:1;font-size:0.68rem" }, `(${c.row},${c.col})`),
                c.weapon !== "n" && m("span.unit-weapon", WEAPON_NAMES[c.weapon]),
                m("span.unit-life", `${c.life}hp`),
              ])
            )
        ),
      ]);
    },
  };
}

function ControlsBar(): m.Component {
  return {
    view() {
      if (!state.game) return m("div.controls-bar");
      const rd = state.game.rounds[state.currentRound];

      return m("div.controls-bar", [
        m("button.btn.btn-icon.btn-sm", { onclick: goToStart, title: "Start" }, "⏮"),
        m("button.btn.btn-icon.btn-sm", { onclick: stepBackward, title: "Step back" }, "◀"),
        m("button.btn.btn-icon.btn-sm", {
          onclick: onPlay,
          title: state.playing ? "Pause" : "Play",
          style: "font-size:0.85rem",
        }, state.playing ? "⏸" : "▶"),
        m("button.btn.btn-icon.btn-sm", { onclick: stepForward, title: "Step forward" }, "▶"),
        m("button.btn.btn-icon.btn-sm", { onclick: goToEnd, title: "End" }, "⏭"),

        m("input.round-slider[type=range]", {
          min: 0,
          max: state.game.num_rounds,
          value: state.currentRound,
          oninput: (e: InputEvent) => {
            setRound(parseInt((e.target as HTMLInputElement).value));
          },
        }),

        m("span.round-label", `${state.currentRound} / ${state.game.num_rounds}`),

        m("span.day-badge", { class: rd.day ? "day" : "night" }, rd.day ? "DAY" : "NIGHT"),

        m("label.speed-label", [
          "Speed: ",
          m("input[type=range]", {
            min: 1,
            max: 30,
            value: state.speed,
            style: "width:60px;vertical-align:middle",
            oninput: (e: InputEvent) => { state.speed = parseInt((e.target as HTMLInputElement).value); },
          }),
          ` ${state.speed}x`,
        ]),
      ]);
    },
  };
}

export function Viewer(): m.Component {
  return {
    oncreate() {
      loadSample();
      window.addEventListener("keydown", (e) => {
        if (!state.game) return;
        if (e.code === "Space") { e.preventDefault(); onPlay(); m.redraw(); }
        else if (e.code === "ArrowRight") { stepForward(); m.redraw(); }
        else if (e.code === "ArrowLeft") { stepBackward(); m.redraw(); }
        else if (e.code === "Home") { goToStart(); m.redraw(); }
        else if (e.code === "End") { goToEnd(); m.redraw(); }
      });
    },
    view() {
      return m("div", { style: "display:flex;flex-direction:column;height:100%" }, [
        m("header.top-bar", [
          m("h1.logo", ["ThePurge", m("span", "Viewer")]),
          state.game && m("span", { style: "font-size:0.72rem;color:var(--text-muted)" },
            `Seed ${state.game.seed} · ${state.game.settings.NUM_DAYS} days · ${state.game.names.length} players`),
        ]),
        m("div.main-layout", [
          m(LeftPanel),
          m("div.center-area", [
            m("div.canvas-wrap", {
              oncreate: renderCanvas,
              onupdate: renderCanvas,
            }, [
              m("canvas"),
              !state.game && !state.loading && m("div", {
                style: "position:absolute;inset:0;display:flex;align-items:center;justify-content:center;color:var(--text-muted);font-size:0.85rem",
              }, "Load a game to begin"),
              state.loading && m("div", {
                style: "position:absolute;inset:0;display:flex;align-items:center;justify-content:center;color:var(--accent);font-size:0.85rem",
              }, "Loading game..."),
            ]),
            m(ControlsBar),
          ]),
          m(RightPanel),
        ]),
      ]);
    },
  };
}
