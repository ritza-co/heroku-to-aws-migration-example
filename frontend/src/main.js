import "./assets/main.css";
import axios from "axios";
import { createApp } from "vue";
import App from "./App.vue";

const app = createApp(App);

app.config.globalProperties.axios = axios;
app.mount("#app");
