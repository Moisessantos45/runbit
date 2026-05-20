# RunBit

> **Desktop playground for JavaScript & TypeScript — no project setup needed.**

RunBit is a lightweight desktop application that lets you write, run and test JavaScript and TypeScript code snippets instantly using **Bun** as the runtime engine. Think of it like a multi-tab scratchpad similar to VS Code, but entirely focused on fast, isolated code execution.

---

## ✨ Features

- **Multi-tab workspace** — open as many independent tabs as you need, each one runs its own code in isolation, just like VS Code windows.
- **Real-time output** — results appear in the side panel as soon as your code finishes executing.
- **TypeScript support** — run `.ts` snippets natively without any extra config thanks to Bun.
- **npm package management** — install and remove npm packages from within the app (powered by `bun add` / `bun remove`).
- **No project boilerplate** — forget about `package.json`, `tsconfig`, Webpack or Vite just to test a snippet.
- **Auto-update notifications** — the app checks for new versions in the background and notifies you natively when one is available.

---

## 🛠 Tech Stack

| Layer      | Technology                              |
|------------|-----------------------------------------|
| Desktop shell | [Wails v2](https://wails.io) (Go)   |
| Runtime    | [Bun](https://bun.sh)                   |
| Frontend   | [Vue 3](https://vuejs.org) + TypeScript |
| Editor     | Monaco Editor (`@guolao/vue-monaco-editor`) |
| Styling    | Tailwind CSS v4                         |
| State      | Pinia                                   |

---

## 🚀 Development

```bash
# Install frontend dependencies (uses Bun)
cd frontend && bun install && cd ..

# Run in dev mode (hot reload)
wails dev -tags webkit2_41
```

## 📦 Building

```bash
# Compile production binary
wails build -clean -tags webkit2_41

# Generate .deb installer (Linux)
./build-deb.sh 1.0.0
```

---

## 📋 Requirements

- Go ≥ 1.23
- Bun ≥ 1.0
- Linux: `libwebkit2gtk-4.1-dev`, `gcc`

---

## 📄 License

Non-Commercial License — © 2025 Moisessantos45.
See [LICENSE](./LICENSE) for details.