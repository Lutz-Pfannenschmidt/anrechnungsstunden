import { defineConfig } from "eslint/config";
import unusedImports from "eslint-plugin-unused-imports";
import eslintPluginAstro from 'eslint-plugin-astro';

export default defineConfig([
    ...eslintPluginAstro.configs.recommended,
    {
        plugins: {
            "unused-imports": unusedImports,
        },
        rules: {
            "no-unused-vars": "off", // or "@typescript-eslint/no-unused-vars": "off",
            "unused-imports/no-unused-imports": "error",
            "unused-imports/no-unused-vars": [
                "warn",
                {
                    vars: "all",
                    varsIgnorePattern: "^_",
                    args: "after-used",
                    argsIgnorePattern: "^_",
                },
            ],
            "sort-imports": [
                "error",
                {
                    ignoreDeclarationSort: true,
                },
            ],
        },
    },
]);
