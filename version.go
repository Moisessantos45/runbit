package main

// Version es inyectada en tiempo de compilación con:
//   -ldflags "-X main.Version=x.y.z"
// Si no se inyecta (ej: wails dev), queda como "dev".
var Version = "dev"
