import vue from 'eslint-plugin-vue'
import vuetify from 'eslint-plugin-vuetify'
import eslintConfigPrettier from "eslint-config-prettier/flat";
import typescriptEslint from 'typescript-eslint'

export default [
    ...vue.configs['flat/base'],
    ...vuetify.configs['flat/base'],
    eslintConfigPrettier,
    {
        files: ['*.vue', '**/*.vue'],
        languageOptions: {parserOptions: {parser: typescriptEslint.parser}},
    },
]
