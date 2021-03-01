module.exports = {
  extends: ['plugin:cypress/recommended',
    'everywhere',
    'everywhere/react',
    'everywhere/typescript',
    'everywhere/jest',
  ],
  parserOptions: {
    parser: '@typescript-eslint/parser',
  },
  rules: {
    "import/extensions": [0]
  },
}
