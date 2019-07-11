shapley_plot <- function(
                         mes_sirets,
                         my_data,
                         model,
                         batch,
                         dir_out = rprojroot::find_rstudio_root_file(
                           "..",
                           "output",
                           "shapley",
                           batch
                         )) {
  assert_that(inherits(my_data, "data.frame"))
  max_periode <- max(my_data$periode)

  h2o::h2o.no_progress()
  pred <- function(model, newdata) {
    results <- as.data.frame(h2o::h2o.predict(model, h2o::as.h2o(newdata)))
    return(results[[3L]])
  }


  x_medium <- c(
    "montant_part_patronale",
    "ratio_dette",
    "ratio_dette_moy12m",
    "etat_proc_collective_num",
    "TargetEncode_code_ape_niveau3",
    "cotisation_moy12m",
    "frais_financier_distrib_APE1",
    "taux_marge_distrib_APE1",
    "montant_part_patronale_past_3",
    "ratio_liquidite_reduite_distrib_APE1",
    "dette_fiscale",
    "ratio_delai_client_distrib_APE1",
    "montant_part_patronale_past_1",
    "montant_part_patronale_past_2",
    "ratio_rend_capitaux_propres",
    "taux_marge",
    "poids_frng_distrib_APE1",
    "delai_fournisseur_distrib_APE1",
    "taux_rotation_stocks_distrib_APE1",
    "effectif",
    "ratio_rentabilite_nette_distrib_APE1",
    "ratio_export_distrib_APE1",
    "TargetEncode_code_ape_niveau2",
    "effectif_past_12",
    "montant_part_ouvriere",
    "financier_court_terme_distrib_APE1",
    # "effectif_entreprise",
    "age",
    "ratio_liquidite_reduite",
    "ratio_productivite_distrib_APE1",
    "frais_financier",
    "financier_court_terme",
    "ratio_delai_client",
    "TargetEncode_code_naf",
    "benefice_ou_perte",
    "taux_rotation_stocks",
    "nombre_etab_secondaire",
    "nbr_etablissements_connus",
    "CA",
    "chiffre_affaires_net_lie_aux_exportations",
    "ratio_dette_delai",
    "ratio_marge_operationnelle_distrib_APE1",
    "poids_frng"
  )



  x_medium_names <- c(
    "Montant part patronale",
    "Ratio dette / cotisation",
    "Moyenne dette/cotisation (12 mois)",
    "Procedure collective en cours",
    "Taux de defaillance dans le secteur d'activite (code APE 3)",
    "Cotisations URSSAF",
    "Comparaison des frais financiers par code NAF",
    "Comparaison du taux de marge par code NAF",
    "Montant part patronale 3 mois en arriere",
    "Comparaison des liquidites reduites par code NAF",
    "Dette fiscale et sociale",
    "Comparaison du delai client par code NAF",
    "Montant part patronale 1 mois en arriere",
    "Montant part patronale 2 mois en arriere",
    "Rendement des capitaux propres",
    "Taux de marge",
    "Poids du frng",
    "Comparaison du delai fournisseur par code NAF",
    "Comparaison du taux de rotation des stocks par code NAF",
    "Effectif salarie",
    "Comparaison de la rentabilite nette par code NAF",
    "ratio_export_distrib_APE1",
    "Taux de defaillance dans le secteur d'activite (code APE 2)",
    "Variation mensuelle d'effectif moyenne sur 12 mois",
    "Montant de la part ouvriere",
    "Comparaison financier court terme par code NAF",
    # "effectif_entreprise",
    "Age de l'entreprise",
    "Ratio des liquidites reduites",
    "Comparaison de la productivite par code NAF",
    "Frais financiers",
    "Frais financiers court terme",
    "Delai client",
    "Taux de defillance dans le secteur d'activite (code NAF)",
    "Resultat net consolide",
    "Taux de rotation des stocks",
    "Nombre d'etablissements secondaires",
    "Nombre d'etablissements connus",
    "Chiffre d'affaire",
    "Chiffre d'affaire net lie aux exportations",
    "Decroissance de la dette pendant un delai URSSAF",
    "Comparaison de la marge operationnelle par code NAF",
    "Poids du frng"
  )

  names(x_medium_names) <- x_medium

  features <- my_data[, x_medium]

  response <- pred(model, features)

  predictor.xgb <- iml::Predictor$new(
    model = model,
    data = features,
    y = response,
    predict.fun = pred,
    class = "classification"
  )

  for (i in seq_along(mes_sirets)) {
    etablissement <- my_data %>%
      filter(siret == mes_sirets[i]) %>%
      filter(periode == max_periode)
    etablissement <- etablissement[, x_medium]
    shap.xgb <- iml::Shapley$new(predictor.xgb, x.interest = etablissement)

    shap_plot <- shap.xgb %>%
      plot()

    shap_plot$data <- shap_plot$data %>%
      mutate(category = unname(x_medium_names[feature])) %>%
      group_by(category) %>%
      summarize(
        feature = unique(category),
        phi = sum(phi),
        phi.var = sum(phi.var)
      ) %>%
      mutate(feature.value = as.factor(category))

    thresh <- 5e-3
    to_remove <- abs(shap_plot$data[, "phi"]) < thresh
    shap_plot$data <- shap_plot$data[!to_remove, ]

    # labels
    shap_plot$labels$x <- ""
    dir.create(dir_out, showWarnings = FALSE)
    ggsave(
      filename = paste0(mes_sirets[i], "_", max_periode, ".png"),
      plot = shap_plot,
      path = dir_out
    )
  }
  h2o::h2o.show_progress()
  return(shap_plot)
}
