import tailwindcss from "@tailwindcss/vite";
// @ts-check
import { defineConfig } from "astro/config";

import vue from "@astrojs/vue";

export default defineConfig({
    vite: {
        plugins: [tailwindcss()],
    },

    integrations: [vue()],
});
