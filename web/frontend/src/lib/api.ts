import type { GameData } from "./types";

const BASE = import.meta.env.VITE_API_URL ?? "";

async function get<T>(url: string): Promise<T> {
  const r = await fetch(`${BASE}${url}`);
  if (!r.ok) throw new Error(await r.text());
  return r.json();
}

async function post<T>(url: string, body: unknown): Promise<T> {
  const r = await fetch(`${BASE}${url}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!r.ok) throw new Error(await r.text());
  return r.json();
}

export const getStatus = () => get<{ status: string; engine_available: boolean }>("/api/status");
export const getSample = () => get<GameData>("/api/sample");
export const getAIs = () => get<string[]>("/api/ais");
export const generateGame = (seed?: number, players?: string[]) =>
  post<GameData>("/api/generate", { seed, players });
export const parseReplay = (replay_text: string) => post<GameData>("/api/parse", { replay_text });
