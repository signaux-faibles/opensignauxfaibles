Mettre-à-jour régulièrement Node.js et les dépendances du projet [`opensignauxfaibles`](https://github.com/signaux-faibles/opensignauxfaibles) (~ 1 fois par mois)

## Mise à jour de Node.js

0. Assurez-vous que `nvm` est installé ([https://github.com/nvm-sh/nvm#installing-and-updating](https://github.com/nvm-sh/nvm#installing-and-updating))
1. Vérifier la dernière version LTS de Node.js sur le [site officiel](https://nodejs.org/)
2. Saisir le numéro de version dans le fichier `js/.nvmrc`
3. Depuis le répertoire `js`, lancer `nvm use` pour basculer sur la nouvelle version (une installation peut être requise, suivre les instructions). 
4. Exécuter tous les tests automatisés (`make test`), pour être sur que tout fonctionne
5. Ne pas oublier de spécifier également ce numéro de version dans la propriété `node-version` dans la configuration du workflow d'intégration continue (`ci.yml`)

## Mise-à-jour des dépendances Node.js

Depuis le répertoire `js`:

0. Usage de la version prévue de Node.js: `nvm use`
1. Mise-à-jour des versions des dépendances dans `package.json`: 
`npx npm-check-updates -u`
2. Installation de ces nouvelles versions de dépendances:
`npm install`
3. Exécution des tests: `make test`
4. Si tout est bon, faire un commit avec `package.json` et `package-lock.json`
