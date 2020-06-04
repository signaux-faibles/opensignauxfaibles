function compte(compte) {
  "use strict";
  const c = f.iterable(compte)
  return (c.length>0)?c[c.length-1]:undefined
}

exports.compte = compte
