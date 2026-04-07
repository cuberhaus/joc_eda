import m from "mithril";
import { Viewer } from "./game/viewer";

(window as any).__mithril = m;

m.mount(document.getElementById("app")!, Viewer);
