{
  "extension": [
    ".js",
    ".ts"
  ],
  "include": [
    "**/*.ts",
    "**/*.js"
  ],
  "exclude": [
    "**/*.d.ts",
    "**/*_test*.ts",
    "**/*_test*.js"
  ],
  "require": [
    "ts-node/register"
  ],
  "cache": true,
  "sourceMap": true, /* nécessaire pour que les instructions couvertes soient associées aux bons numéros de lignes des fichiers TypeScript avant transpilation, cf https://github.com/signaux-faibles/opensignauxfaibles/pull/348#issuecomment-822462720 */
  "instrument": true,
  "all": true
}
