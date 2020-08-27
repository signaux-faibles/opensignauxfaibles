module.exports = {
  parser: "@typescript-eslint/parser",
  ignorePatterns: ["node_modules/"],
  extends: [
    "eslint:recommended",
    "plugin:prettier/recommended",
    "plugin:@typescript-eslint/eslint-recommended",
    "plugin:@typescript-eslint/recommended",
    "prettier/@typescript-eslint",
  ],
  plugins: ["prettier", "@typescript-eslint"],
  parserOptions: {
    tsconfigRootDir: __dirname,
    ecmaVersion: 2019,
  },
  rules: {
    "prettier/prettier": "error", // tout problème de formatage detecté par prettier sera reporté comme une erreur par `$ npm run lint`
    "@typescript-eslint/camelcase": 0, // tolérer l'usage de noms en snake case (avec underscores)
    eqeqeq: ["warn", "always"], // pour encourager l'usage de opérateurs d'égalité stricts (=== au lieu de ==, et !== au lieu de !=)
    "object-shorthand": ["warn", "always"], // pour s'aligner avec codacy
  },
  env: {
    node: true, // permet l'usage du global "module"
  },
}
