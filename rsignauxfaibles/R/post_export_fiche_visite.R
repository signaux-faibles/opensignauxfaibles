#' Export d'une fiche visite
#'
#' @param sirets Vecteur de sirets des établissements dont il faut exporter la
#' fiche
#' @param database Base de donnée Mongodb à requêter
#' @param collection Collection mongodb à requêter
#' @param batch Choix du batch pour lequel se fait l'export
#' @param with_urssaf. logical. Si TRUE, les informations du montant des
#'  dettes URSSAF sont incluses à la fiche.
#' @param folder Nom du dossier dans lequel la fiche est exportée. Le dossier
#'  doit préexister. TODO: chemin d'accès écrit en dur; à corriger !
#'
#' @return NULL
#' @export
#'
#' @examples
export_fiche_visite <- function(
  sirets,
  database = "test_signauxfaibles",
  collection = "Features",
  batch,
  with_urssaf = FALSE,
  folder = batch){

  for (i in seq_along(sirets)) {

    elementary_info  <- connect_to_database(
      database = database,
      collection = collection,
      batch = batch,
      siren = substr(sirets[i], 1, 9),
      fields = c("siret", "raison_sociale", "periode", "code_ape_niveau3"),
      min_effectif = 0
    ) %>%
      filter(siret == sirets[i]) %>%
      head(n = 1)

    raison_sociale <- elementary_info %>%
      .$raison_sociale

    rmarkdown::render("/home/pierre/Documents/opensignauxfaibles/rsignauxfaibles/R/post_fiche_visite.Rmd",
      params  = list(
        siret = sirets[i],
        batch = batch,
        database = database,
        collection = collection,
        with_urssaf = with_urssaf
      ),
      output_file = paste0("/home/pierre/Documents/opensignauxfaibles/output/Fiches/", folder, "/Fiche_visite_", stringr::str_replace_all(raison_sociale, "[^[:alnum:]]", "_"), ".pdf"),
      clean = TRUE
    )
    }
}
