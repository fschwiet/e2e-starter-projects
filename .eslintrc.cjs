module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  plugins: ['@typescript-eslint'],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'prettier',
  ],
  overrides: [
    {
      files: ['src/main/**/*', 'tests/unit/**/*', 'tests/e2e/**/*'],
      env: { node: true },
    },
    {
      files: ['src/renderer/**/*'],
      env: { browser: true },
    },
    {
      files: ['*.cjs'],
      env: { node: true },
      rules: {
        '@typescript-eslint/no-require-imports': 'off',
      },
    },
  ],
};
