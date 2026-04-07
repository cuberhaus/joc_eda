import type { RoundData, Citizen, Barricade, GameData } from "../lib/types";

const PLAYER_COLORS = ["#22c55e", "#ef4444", "#3b82f6", "#a855f7"];
const BUILDING_COLOR = "#2a2a3e";
const BUILDING_EDGE = "#3a3a50";
const STREET_DAY = "#1e1e30";
const STREET_NIGHT = "#14141f";
const FOOD_COLOR = "#f59e0b";
const MONEY_COLOR = "#fbbf24";
const GUN_COLOR = "#94a3b8";
const BAZOOKA_COLOR = "#f97316";

export class GameCanvas {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;
  private tileSize = 0;
  private offsetX = 0;
  private offsetY = 0;

  constructor(canvas: HTMLCanvasElement) {
    this.canvas = canvas;
    this.ctx = canvas.getContext("2d")!;
  }

  render(game: GameData, round: number, animProgress: number) {
    const rd = game.rounds[round];
    const nextRd = round < game.num_rounds ? game.rounds[round + 1] : null;
    const { rows, cols } = game;

    const rect = this.canvas.parentElement!.getBoundingClientRect();
    const dpr = window.devicePixelRatio || 1;
    this.canvas.width = rect.width * dpr;
    this.canvas.height = rect.height * dpr;
    this.canvas.style.width = `${rect.width}px`;
    this.canvas.style.height = `${rect.height}px`;

    const ctx = this.ctx;
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

    const maxTileW = rect.width / cols;
    const maxTileH = rect.height / rows;
    this.tileSize = Math.floor(Math.min(maxTileW, maxTileH));
    if (this.tileSize < 4) this.tileSize = 4;

    this.offsetX = Math.floor((rect.width - this.tileSize * cols) / 2);
    this.offsetY = Math.floor((rect.height - this.tileSize * rows) / 2);

    ctx.fillStyle = "#0a0a12";
    ctx.fillRect(0, 0, rect.width, rect.height);
    ctx.save();
    ctx.translate(this.offsetX, this.offsetY);

    const isDay = rd.day === 1;
    const streetColor = isDay ? STREET_DAY : STREET_NIGHT;

    for (let r = 0; r < rows; r++) {
      for (let c = 0; c < cols; c++) {
        const cell = rd.grid[r][c];
        const x = c * this.tileSize;
        const y = r * this.tileSize;
        const s = this.tileSize;

        if (cell === "B") {
          ctx.fillStyle = BUILDING_COLOR;
          ctx.fillRect(x, y, s, s);
          ctx.strokeStyle = BUILDING_EDGE;
          ctx.lineWidth = 0.5;
          ctx.strokeRect(x + 0.5, y + 0.5, s - 1, s - 1);
        } else {
          ctx.fillStyle = streetColor;
          ctx.fillRect(x, y, s, s);
          ctx.strokeStyle = isDay ? "#252538" : "#1a1a28";
          ctx.lineWidth = 0.3;
          ctx.strokeRect(x, y, s, s);
        }
      }
    }

    this.drawBarricades(rd.barricades);
    this.drawItems(rd, nextRd, animProgress);
    this.drawCitizens(rd, nextRd, animProgress);

    ctx.restore();
  }

  private drawBarricades(barricades: Barricade[]) {
    const ctx = this.ctx;
    const s = this.tileSize;
    for (const bar of barricades) {
      const x = bar.col * s;
      const y = bar.row * s;
      const strength = Math.min(4, Math.ceil(bar.resistance / 80));
      const alpha = 0.3 + strength * 0.15;
      ctx.fillStyle = PLAYER_COLORS[bar.player] + Math.round(alpha * 255).toString(16).padStart(2, "0");
      ctx.fillRect(x + 1, y + 1, s - 2, s - 2);
      ctx.strokeStyle = PLAYER_COLORS[bar.player];
      ctx.lineWidth = 1;
      ctx.setLineDash([2, 2]);
      ctx.strokeRect(x + 1.5, y + 1.5, s - 3, s - 3);
      ctx.setLineDash([]);
    }
  }

  private drawItems(rd: RoundData, nextRd: RoundData | null, anim: number) {
    const ctx = this.ctx;
    const s = this.tileSize;
    const half = s / 2;
    const itemR = Math.max(2, s * 0.22);

    for (let r = 0; r < rd.grid.length; r++) {
      for (let c = 0; c < rd.grid[r].length; c++) {
        const cell = rd.grid[r][c];
        const cx = c * s + half;
        const cy = r * s + half;
        const visible = !nextRd || nextRd.grid[r][c] === cell || anim < 0.5;

        if (cell === "F" && visible) {
          ctx.fillStyle = FOOD_COLOR;
          ctx.beginPath();
          ctx.arc(cx, cy, itemR, 0, Math.PI * 2);
          ctx.fill();
        } else if (cell === "M" && visible) {
          ctx.fillStyle = MONEY_COLOR;
          const ir = itemR * 0.9;
          ctx.fillRect(cx - ir, cy - ir, ir * 2, ir * 2);
        } else if (cell === "G" && visible) {
          ctx.fillStyle = GUN_COLOR;
          ctx.beginPath();
          ctx.moveTo(cx, cy - itemR);
          ctx.lineTo(cx + itemR, cy + itemR);
          ctx.lineTo(cx - itemR, cy + itemR);
          ctx.closePath();
          ctx.fill();
        } else if (cell === "Z" && visible) {
          ctx.fillStyle = BAZOOKA_COLOR;
          ctx.beginPath();
          const ir = itemR * 1.1;
          ctx.moveTo(cx, cy - ir);
          ctx.lineTo(cx + ir * 0.7, cy - ir * 0.3);
          ctx.lineTo(cx + ir, cy + ir);
          ctx.lineTo(cx, cy + ir * 0.5);
          ctx.lineTo(cx - ir, cy + ir);
          ctx.lineTo(cx - ir * 0.7, cy - ir * 0.3);
          ctx.closePath();
          ctx.fill();
        }
      }
    }
  }

  private drawCitizens(rd: RoundData, nextRd: RoundData | null, anim: number) {
    const ctx = this.ctx;
    const s = this.tileSize;
    const citizenMap = new Map<number, Citizen>();
    if (nextRd) nextRd.citizens.forEach(c => citizenMap.set(c.id, c));

    for (const cit of rd.citizens) {
      let dx = cit.col * s;
      let dy = cit.row * s;

      if (nextRd && anim > 0) {
        const next = citizenMap.get(cit.id);
        if (next) {
          dx = cit.col * s + (next.col - cit.col) * s * anim;
          dy = cit.row * s + (next.row - cit.row) * s * anim;
        }
      }

      const cx = dx + s / 2;
      const cy = dy + s / 2;
      const color = PLAYER_COLORS[cit.player] || "#888";
      const unitR = Math.max(3, s * 0.32);

      if (cit.type === "b") {
        ctx.fillStyle = color;
        ctx.fillRect(cx - unitR, cy - unitR, unitR * 2, unitR * 2);
        ctx.fillStyle = "#0f0f1a";
        ctx.fillRect(cx - unitR * 0.35, cy - unitR * 0.35, unitR * 0.7, unitR * 0.7);
      } else {
        ctx.fillStyle = color;
        ctx.beginPath();
        ctx.arc(cx, cy, unitR, 0, Math.PI * 2);
        ctx.fill();

        if (cit.weapon !== "n" && cit.weapon !== "h") {
          ctx.fillStyle = cit.weapon === "b" ? BAZOOKA_COLOR : GUN_COLOR;
          ctx.beginPath();
          ctx.arc(cx + unitR * 0.5, cy - unitR * 0.5, Math.max(1.5, unitR * 0.3), 0, Math.PI * 2);
          ctx.fill();
        }
      }

      if (s >= 14) {
        const maxLife = cit.type === "b" ? 60 : 100;
        const lifePct = Math.min(1, cit.life / maxLife);
        const barW = unitR * 2;
        const barH = Math.max(1.5, s * 0.06);
        const barY = cy + unitR + 2;
        ctx.fillStyle = "#0f0f1a80";
        ctx.fillRect(cx - barW / 2, barY, barW, barH);
        ctx.fillStyle = lifePct > 0.5 ? "#22c55e" : lifePct > 0.25 ? "#f59e0b" : "#ef4444";
        ctx.fillRect(cx - barW / 2, barY, barW * lifePct, barH);
      }
    }
  }
}
