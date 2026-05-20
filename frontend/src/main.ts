import { createApp } from "vue";
import { install as VueMonacoEditorPlugin } from "@guolao/vue-monaco-editor";
import { createPinia } from "pinia";
import App from "./App.vue";
import "./style.css";

const app = createApp(App);
app.use(createPinia());

app.use(VueMonacoEditorPlugin, {
  paths: {
    vs: "https://cdn.jsdelivr.net/npm/monaco-editor@0.52.2/min/vs",
  },
});

app.mount("#app");
