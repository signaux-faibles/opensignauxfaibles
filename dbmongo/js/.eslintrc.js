module.exports = {
  parser: "@typescript-eslint/parser",
  ignorePatterns: ["lib/", "node_modules/"],
  extends: [
    "eslint:recommended",
    "plugin:prettier/recommended",
    "plugin:@typescript-eslint/eslint-recommended",
    "plugin:@typescript-eslint/recommended",
    "prettier/@typescript-eslint",
  ],
  plugins: ["prettier", "@typescript-eslint"],
  parserOptions: {
    project: "./tsconfig.json",
    ecmaVersion: 2019,
  },
  rules: {
    "prettier/prettier": "error",
    "@typescript-eslint/camelcase": 0,
  },
  env: {
    node: true,
  },
}
