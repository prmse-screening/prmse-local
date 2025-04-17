import vue from 'eslint-plugin-vue'
import vuetify from 'eslint-plugin-vuetify'
import typescriptEslint from 'typescript-eslint'

export default [
    ...vue.configs['flat/base'],
    ...vuetify.configs['flat/base'],
    {
        files: ['**/*.vue'],
        languageOptions: { parserOptions: { parser: typescriptEslint.parser } },
    },
]
