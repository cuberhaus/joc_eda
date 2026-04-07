export interface Citizen {
  type: "b" | "w";
  id: number;
  player: number;
  row: number;
  col: number;
  weapon: "n" | "h" | "g" | "b";
  life: number;
}

export interface Barricade {
  player: number;
  row: number;
  col: number;
  resistance: number;
}

export interface Command {
  id: number;
  action: "m" | "b";
  dir: "u" | "d" | "l" | "r";
}

export interface RoundData {
  grid: string[];
  citizens: Citizen[];
  barricades: Barricade[];
  round: number;
  day: number;
  scores: number[];
  cpu: string[];
  commands: Command[];
}

export interface GameData {
  sec_game: boolean;
  seed: number;
  version: string;
  settings: Record<string, number>;
  names: string[];
  num_rounds: number;
  rows: number;
  cols: number;
  rounds: RoundData[];
}
