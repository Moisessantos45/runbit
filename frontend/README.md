# RunBit — Frontend

Vue 3 + TypeScript frontend for the RunBit desktop app, bundled with Vite and powered by Bun.

---

## Tech Stack

| Package | Version | Purpose |
|---------|---------|---------|
| `vue` | ^3.5.34 | UI framework |
| `@guolao/vue-monaco-editor` | ^1.6.0 | Code editor component (Monaco / VS Code engine) |
| `pinia` | ^3.0.4 | State management |
| `tailwindcss` | ^4.3.0 | Utility-first CSS framework |
| `@tailwindcss/vite` | ^4.3.0 | Tailwind Vite plugin |

### Dev Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `vite` | ^8.0.12 | Build tool & dev server |
| `@vitejs/plugin-vue` | ^6.0.6 | Vue support for Vite |
| `typescript` | ~6.0.2 | TypeScript compiler |
| `vue-tsc` | ^3.2.8 | Type-check for Vue SFCs |
| `@vue/tsconfig` | ^0.9.1 | Shared TS config for Vue |
| `@types/node` | ^24.12.3 | Node type definitions |
| `@babel/types` | ^7.18.10 | Required by Monaco editor |

---

## Notes

- The frontend communicates with the Go backend via the **Wails runtime** (`window.runtime` / generated bindings in `wailsjs/`).
- Events like `runner:stdout`, `runner:stderr`, `runner:done`, `runner:ready`, and `app:update-available` are emitted from Go and consumed here.
- Do **not** run `npm install` — this project uses **Bun** exclusively.
