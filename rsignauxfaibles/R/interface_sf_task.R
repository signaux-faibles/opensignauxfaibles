#' Documentation des informations de connection à mongodb
#'
#' @param url `character(1)` \cr url to the database in mongodb uri format.
#' @param database `character(1)` \cr Nom de la base de données vers laquelle
#'   param exporter. Par défaut, celle stockée dans \code{task}.
#' @param collection `character(1)` \cr Nom de la collection vers laquelle
#'   exporter. Par défaut, celle stockée dans \code{task}.
#' @name mongodb_connection
NULL

#' Documentation tâche générique
#'
#' @param task `[sf_task]` \cr Objet s3 de type sf_task
#' @param tracker `[mllogr::Tracker]` \cr Logger of ML experiment
#' @name generic_task
NULL

#' Initialiser une tâche d'apprentissage
#'
#' Un objet s3 de type sf_task est défini, dans lequel seront défini et
#' stockés les tâches intermédiaires et les résultats de l'apprentissage.
#'
#' Si `tracker` est `NULL`, alors un nouveau Tracker est créé
#'
#' @param verbose `logical(1)` \cr
#'   Active ou désactive le log des actions ultérieures.
#' @inheritParams mongodb_connection
#' @inheritParams generic_task
#' @param experiment_name `character()` \cr
#'   Quel est l'objet de cette tâche ?
#' @param experiment_description `character()` \cr
#'   Descriptions supplémentaires sur l'expérimentation en cours.
#'
#' @return `[rsignauxfaibles::sf_task]` \cr
#'   Un objet sf_task avec un attribut de type `logical` "verbose", qui définit le niveau de log,
#'   ainsi qu'un attribut "to_log" de type `list` dans lequel seront stockés des
#'   informations spécifiques pour le log.
#'
#'  @export
sf_task <- function(
  verbose,
  database = "test_signauxfaibles",
  collection = "Features",
  experiment_name,
  experiment_description,
  tracker = NULL
  ){

  res <- list(database = database, collection = collection)
  class(res) <- "sf_task"
  attr(res, "verbose") <- verbose
  if (is.null(tracker) && require(MLlogr)){

    res[["tracker"]] <- MLlogr::Tracker$new(
      #TODO: Mieux gérer la base de logging
      # f4737679-916d-4aeb-84ba-958046f4ca31
      database = database,
      collection = collection,
      control = list(id_columns = c("siret", "periode"))
    )
    res[["tracker"]]$set(
      experiment_name = experiment_name,
      experiment_description = experiment_description
    )
  }
  return(res)
}

#' Active ou désactive le logging
#'
#' Prend en compte l'attribut "verbose" de l'objet task pour fixer le bon
#' niveau de logging
#' @param task `[sf_task]` \cr Objet s3 de type sf_task
#' @return TRUE
set_verbose_level <- function(task){
  require(logger)
  if (attr(task, "verbose")){
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }
  return(TRUE)
}

#' Vérification de champs
#'
#' Vérifie si les champs qui vont être écrits sont déjà existant, et le cas
#' échéant vont être écrasés.
#' @param task `[sf_task]` \cr Objet s3 de type sf_task
#' @param field_names `character()` \cr Nom des champs à vérifier.
#' @return Nom des champs écrasés, `character(0)` sinon.
check_overwrites <- function(task, field_names){
  set_verbose_level(task)
  overwrite <- intersect(field_names, names(task))
  if (length(overwrite) > 1){
    log_info('Les champs {paste(overwrite, collapse = ",")} sont ecrases avec
      les nouvelles valeurs.')
  }
  return(overwrite)
}

#' Chargement de données historiques
#'
#' Charge les données historiques de signaux faibles et les stocke dans un
#' champ "hist_data" de l'objet \code{task}
#'
#' @param batch `character(1)` \cr Batch auquel doit être importées les
#'   données. Les modifications opérées par les batchs ultérieurs sont
#'   ignorées.
#' @inheritParams mongodb_connection
#' @param subsample `integer(1)` \cr Nombre d'objets (c'est-à-dire de couples
#'   siret x periode) à échantillonner.
#' @param fields `character()` \cr Noms des champs à requêter dans la base de
#'   données. Doit contenir "siret" et "periode". Si égal à \code{NULL}, alors
#'   charge tous les champs disponibles.
#' @param date_inf `Date` \cr Date inférieure du chargement des données.
#' @param date_sup `Date` \cr Date supérieure (exclue) du chargement des
#'   données.
#' @param min_effectif `integer(1)` \cr Limite basse du filtrage de l'effectif
#'   (la limite est incluse)
#' @param siren `character()` \cr Liste de sirens à exporter. Si égale à
#'   \code{NULL}, charge tous les sirens disponibles.
#' @param code_ape `character()` \cr Liste de codes APE à exporter. Si égale à
#'   \code{NULL}, charge tous les codes disponibles.
#'
#' @return `[sf_task]` \cr
#'   L'objet \code{task} donné en entrée auquel le champs "hist_data" a été
#'   ajouté (ou écrasé), contenant un  data.frame() avec les colonnes incluses dans le paramètre d'entrée
#'  \code{fields}, et pour chaque ligne un couple unique siret x periode.
#'
#' @describeIn load_hist_data
#'
#' @export
load_hist_data.sf_task <- function(
  task,
  batch,
  database = task[["database"]],
  collection = task[["collection"]],
  subsample = 200000L,
  fields = get_fields(training = FALSE),
  date_inf = as.Date("2015-01-01"),
  date_sup = as.Date("2017-01-01"),
  min_effectif = 10L,
  siren = NULL,
  code_ape = NULL
  ){
  set_verbose_level(task)

  log_info("Chargement des donnees historiques.")

  hist_data <- connect_to_database(
    database,
    collection,
    batch,
    min_effectif = min_effectif,
    siren = NULL,
    date_inf = date_inf,
    date_sup = date_sup,
    fields = fields,
    code_ape = NULL,
    subsample = subsample,
    verbose = attr(task, "verbose")
    )

  if (nrow(hist_data) > 1) {
    log_info("Les donnees ont ete chargees avec succes.")
  } else {
    log_warn("Aucune donnee n'a ete chargee. Veuillez verifier la requete.")
  }
  check_overwrites(task, "hist_data")
  task[["hist_data"]] <- hist_data
  return(task)
}

#' Chargement de nouvelles données
#'
#' @param task `[sf_task]` \cr Objet s3 de type sf_task
#' @param periods `[Date()]` \cr Périodes d'intérêt, auquels charger les
#'   données. Des périodes supplémentairs peuvent être chargées selon la
#'   valeur de rollback_months.
#' @inheritParams connect_to_database
#' @param rollback_months `integer(1)`\cr Nombre de mois précédant le premier mois de
#'   `periods` à charger. Permet d'effectuer des calculs de différences ou de
#'   moyennes glissantes pour les périodes d'intérêt.
#'
#' @describeIn load_new_data
#'
#' @return `[sf_task]` \cr
#'   L'objet \code{task} donné en entrée auquel le champs "new_data" a été
#'   ajouté (ou écrasé), contenant un  data.frame() avec les colonnes incluses dans le paramètre d'entrée
#'  \code{fields}, et pour chaque ligne un couple unique siret x periode.
#' @export
load_new_data.sf_task <- function(
  task,
  periods,
  batch,
  database = task[["database"]],
  collection = task[["collection"]],
  fields = get_fields(training = FALSE),
  min_effectif = 10L,
  rollback_months = 1L
  ){

  set_verbose_level(task)

  log_info("Loading data from last batch")
  task[["new_data"]] <- connect_to_database(
    database = database,
    collection = collection,
    batch = batch,
    date_inf = min(periods) %m-% months(rollback_months),
    date_sup = max(periods) %m+% months(1),
    min_effectif = min_effectif,
    fields = fields
  )

  if ("periode" %in% fields &&
    max(task[["new_data"]]$periode) != max(periods)) {
    log_warn("Data is missing at actual period !")
  }
  return(task)
}


#' Scission des données en échantillon d'entraînement, de validation et de
#' test.
#'
#' Scinde les données historiques en échantillon d'entraînement, de validation et de
#' test, selon les proportions souhaitées. S'assure que deux établissements de la même entreprise ne soient pas
#' à la fois dans deux échantillons différents pour éviter la fuite
#' d'information d'un échantillon vers l'autre.
#'
#'  La fraction de l'échantillon de test est calculée par
#'  1 - frac_train - frac_val. (frac_train + frac_val) doit donc être inférieur
#'  à 1. Le seul cas où cette condition n'est pas testée est lorsque frac_train
#'  = 1.
#'
#' @inheritParams generic_task
#' @inheritParams split_snapshot_rdm_month
#' @describeIn split_data
#'
#' @return `[sf_task]` \cr
#'   L'objet \code{task} donné en entrée auquel les champs "train_data",
#'   "validation_data" et "test_data" ont été
#'   ajoutés (ou écrasés), chacun contenant un data.frame() avec les colonnes
#'   de `task[["hist_data"]]` et un sous-ensemble (possiblement vide) de ces
#'   lignes.
#'
#' @export
split_data.sf_task <- function(
  task,
  fracs = c(0.6, 0.2, 0.2),
  names = c("train", "validation", "test"),
  tracker = task[["tracker"]]
  ){

  set_verbose_level(task)

  log_info("Les donnees historiques sont scindes en echantillons
    d'entrainement, de test et de validation")

  assertthat::assert_that("hist_data" %in% names(task),
    msg = "Please load historical data before holding out test data")

  if ((length(fracs) == 1 && fracs == 1) || identical(fracs,c(1,0,0))) {
    task[["train_data"]] <- task[["hist_data"]]
  } else {

    res <- split_snapshot_rdm_month(
      data_sample = task[["hist_data"]],
      fracs = fracs,
      names = names
      )

    for (name in names){
      task[[paste0(name, "_data")]] <- task[["hist_data"]] %>%
        semi_join(res[[name]], by = c("siret", "periode"))
    }
  }


  if (!is.null(tracker)){
    names(fracs) <- names
    tracker$set(resampling_strategy = "holdout", train_val_test_shares = fracs)
  }
  return(task)
}

#' Préparation des échantillons
#'
#' Prépare les échantillons souhaités (listés dans `data_names`) pour
#' l'entraînement ou la prédiction de l'algorithme. Cf `[prepare_data]` pour
#' la nature de cette préparation.
#'
#' @inheritParams generic_task
#' @param data_names `character()` \cr Vecteur de noms des données à préparer.
#'   Doivent-être des noms de champs valides de `task`.
#'
#' @return `[sf_task]` \cr
#'   L'objet \code{task} donné en entrée auquel les champs de données
#'   préparées ont été ajoutées, avec une convention de nommage d'apposition
#'   du préfixe "prepared_" aux noms des données noms préparées (par exemple:
#'   "prepared_train_data" correspond aux données de "train_data" préparées).
#'
#' @describeIn prepare
#'
#' @export
prepare.sf_task <- function(
  task,
  data_names = c(
    "train_data",
    "validation_data",
    "test_data",
    "new_data"
    ),
  tracker = task[["tracker"]]
  ){


  ## TODO: Should not fail with empty frames
  # 93c281b7-d567-4026-aba4-e13d289b9170

  set_verbose_level(task)

  aux_function <- function(task, name){
    if (name == "train_data"){
      test_or_train <- "train"
    } else {
      test_or_train <- "test"
      assertthat::assert_that("preparation_map" %in% names(task),
        msg = 'No preparation map has been found. Have you loaded a task, or
        prepared the "train_data" first ?')
    }

    if (name %in% names(task)){

      log_info("Preparing {name} for training and predicting")
      out <-  prepare_frame(
        data_to_prepare = task[[name]],
        save_or_load_map = FALSE,
        te_map = task[["preparation_map"]],
        test_or_train = test_or_train,
        outcome = "outcome"
        )

      task[[paste0("prepared_", name)]]  <- out[["data"]]
      if (name == "train_data"){
        task[["preparation_map"]]  <- out[["te_map"]]
      }

      return(task)

    } else {
      log_warn("There is no {name} in current task")
      return(task)
    }
  }

  for (name in data_names){
    task  <- aux_function(task, name)
  }

  if (!is.null(tracker)){
    tracker$set(preprocessing_strategy = "H2OFrame, Target encoding with H2O")
  }

  return(task)
}

#' Optimize the hyperparameters of a model
#'
#' Optimise les hyperparamètres du modèle. Nécessite des données
#' d'entraînement et de validation préparées.
#'
#' @param task `[sf_task]` \cr Objet s3 de type sf_task
#' @param fields `character()` \cr Liste des variables pour l'entraînement. Cf
#' `[get_fields]` pour les variables par défaut.
#' @param n_init `integer(1)` \cr Nombre d'évaluations aléatoires initiales.
#' @param n_iter `integer(1)` \cr Nombre d'itérations (d'évaluations) d'optimisation.
#' @param train_pipe `function` \cr Fonction d'entraînement et d'évaluation
#'   compatible avec l'optimisation avec "mlrMBO" (donc avec des attributs
#'   spécifiques), notamment comme celles créées avec
#'   `[smoof::makeSingleObjectiveFunction]`
#' @param optim_bounds `ParamSet` \cr Objet précisant l'espace exploré par
#'   l'optimisation, tel que donné par `[ParamHelpers::ParamSet]`. Si NULL,
#'   alors des limites par défaut sont fixées.
#'
#' @describeIn optimize_hyperparameters
#'
#' @return La `task` donnée en entrée, auquel a été ajouté un champ
#' "model_parameters" avec les paramètres optimaux calculées.
#' @export
#'
optimize_hyperparameters.sf_task <- function( #nolint
  task,
  fields = get_fields(training = TRUE),
  n_init = NULL,
  n_iter = 12,
  train_pipe = NULL,
  optim_bounds = NULL
){

  require(DiceKriging) #nolint
  require(rgenoud)
  require(mlrMBO) #nolint
  require(ParamHelpers) #nolint

  if (is.null(optim_bounds)){
    optim_bounds <- ParamHelpers::makeParamSet(
      ParamHelpers::makeNumericParam("learn_rate",  lower = 0.003, upper = 0.2),
      ParamHelpers::makeIntegerParam("max_depth", lower = 2L, upper = 12L),
      ParamHelpers::makeIntegerParam("ntrees", lower = 10L, upper = 300L),
      ParamHelpers::makeIntegerParam(
        "min_child_weight",
        lower = 1L,
        upper = 9L
      )
      )
  }

  if (is.null(train_pipe)){
    #1 Objective function
    train_pipe <- smoof::makeSingleObjectiveFunction(
      name = "lgb_pipe",
      fn = function(
        x
        ){

        new_task <- train(
          task,
          parameters = as.list(x),
          fields = fields
          )
        new_task <- predict(new_task, data_names = c("validation_data"))
        new_task <- evaluate(
          new_task,
          data_name = c("validation_data"),
          plot = FALSE
        )
        # TODO: Bayesian optimization criteria should be more flexible
        # 77d91ced-beac-4e91-9b48-e8fd16e956ee
        aucpr <- new_task[["model_performance"]] %>% .$evaluation %>% .[[1]]
        cat(aucpr)
        return(aucpr)
      },
      par.set = optim_bounds,
      minimize = FALSE
      )
  }

    set.seed(159)
  #2 Initial design
  if (!is.null(n_init)){
    des <- ParamHelpers::generateDesign(n = n_init,
      par.set = ParamHelpers::getParamSet(train_pipe),
      fun = lhs::randomLHS)
  } else {
    des <- NULL
  }

  #3 Surrogate model
  surrogate  <- mlr::makeLearner(
    "regr.km",
    predict.type = "se",
    covtype = "matern3_2",
    control = list(trace = FALSE)
  )

  #4 Control
  control <- mlrMBO::makeMBOControl()
  control <- mlrMBO::setMBOControlTermination(
    control,
    iters = n_iter
  )
  control <- mlrMBO::setMBOControlInfill( #nolint
    control,
    crit = mlrMBO::makeMBOInfillCritEI() #nolint
  )

  #5 optimization
  run <- mlrMBO::mbo(fun = train_pipe,
    design = des,
    control = control,
    show.info = TRUE)

  task[["model_parameters"]] <- as.list(run$x)

  all_res <- run$opt.path$env$path
  if (is.null(n_init)){
    n_init  <- 4 * smoof::getNumberOfParameters(train_pipe)
  }
  run$opt.path$env$path  %>%
    mutate(round = row_number()) %>%
    mutate(type = case_when(
        round <= n_init ~ "initial design",
        TRUE ~ "mlrMBO optimization")) %>%
    ggplot2::ggplot(ggplot2::aes(x = round, y = y, color = type)) +
    ggplot2::geom_point() +
    ggplot2::labs(title = "mlrMBO optimization")

  return(task)
}

#' Entraînement de l'algorithme
#'
#' Entraîne un algorithme sur les données `task[["train_data"]]`. Cf
#' `[train_light_gradient_boosting]` pour plus de détails sur le modèle
#' entraîné.
#'
#' @inheritParams generic_task
#' @param fields `character()` \cr Liste des variables pour l'entraînement. Cf
#' `[get_fields]` pour les variables par défaut.
#' @param parameters `list(any())` \cr Paramètres du modèle. Cf
#' `[train_light_gradient_boosting]`. Si `NULL`, alors des paramètres par
#' défaut sont sélectionnés.
#' @param seed `integer()`\cr Graîne pour que les calculs aléatoires soient
#' reproductibles.
#'
#'
#' @return `[sf_task]` \cr L'objet `task` donné en entré, auquel a été ajouté
#' (ou écrasé) le champs "model", dans lequel est stocké un modèle compatible
#' avec la fonction `[prepare.sf_task]`.
#'
#' @export
train.sf_task <- function(
  task,
  fields = get_fields(training = TRUE),
  parameters = NULL,
  seed = 123,
  tracker = task[["tracker"]]
  ){

  assertthat::assert_that(
    "prepared_train_data" %in% names(task),
    msg = "task does not contain prepared train data."
    )

  if (is.null(parameters) &&
    (!"model_parameters" %in% names(task))){
    parameters  <- list(
      learn_rate = 0.1,
      max_depth = 4,
      ntrees = 60,
      min_child_weight = 1
      )

    task[["model_parameters"]] <- parameters
  } else if (is.null(parameters)){
    parameters <- task[["model_parameters"]]
  }

  assertthat::assert_that(
    is.list(parameters)
    )

  assertthat::assert_that(
    all(c("learn_rate", "max_depth", "ntrees", "min_child_weight") %in%
      names(parameters)),
    msg = paste("Following parameters are missing: ",
      dplyr::setdiff(c("learn_rate", "max_depth", "ntrees", "min_child_weight"),
      names(parameters)),
      sep = ", ")
    )

  set_verbose_level(task)

  log_info("Model is being trained.")

  task[["features"]] <- fields

  model <- train_light_gradient_boosting(
    h2o_train_data = task[["prepared_train_data"]],
    x_fields_model = fields,
    save_results = FALSE,
    learn_rate = parameters[["learn_rate"]],
    max_depth = parameters[["max_depth"]],
    ntrees = parameters[["ntrees"]],
    min_child_weight = parameters[["min_child_weight"]],
    seed = seed
    )

  log_info("Model trained_successfully")
  task[["model"]] <- model

  if (!is.null(tracker)){
    #TODO: Why does it print "model"?
    # e51f22ca-ed0d-434b-be1e-d57547092a22
     tracker$set(
       model_name  = "light gradient boosting",
       model  = "model",
       model_target  = "18 mois, defaut et defaillance"
       )
  }

  return(invisible(task))
}

load.sf_task <- function(task){
  set_verbose_level(task)
  task[["model"]] <- load_h2o_object("lgb", "model", last = TRUE)
  task[["preparation_map"]] <- load_h2o_object("te_map", "temap", last = TRUE)

  return(task)
}

save.sf_task <- function(task, ...) {
  set_verbose_level(task)

  assertthat::assert_that(all(c("model", "preparation_map") %in% names(task)),
    msg = "Task should have a model and a preparation_map to be saved")
  save_h2o_object(task[["model"]], "lgb")
  save_h2o_object(task[["preparation_map"]], "preparation_map")

  return(task)
}

#' Predict from a trained model
#'
#' Predict on some data.
#' Data should be prepared.
#'
#' @param object `sf_task()` \cr Objet `sf_task` avec un modèle entraîné dans
#'   le champs `model`. Les données doivent avoir été préparées.
#' @param data_names `character()` \cr Nom des données sur lesquelles prédire.
#'
#' @return `sf_task()`\cr
#'   Task donnée en entrée, pour laquelle une colonne "score" a été rajoutée
#'   aux données nommées dans `data_names`
#' @export
predict.sf_task <- function(
  object,
  data_names = c(
    "new_data",
    "train_data",
    "validation_data",
    "test_data"
    )
  ){

  task  <- object
  set_verbose_level(task)

  assertthat::assert_that(all(c("model") %in% names(task)),
    msg = "Task should have a model to predict on new
    data")

    predict_on_given_data <- function(data_name, task){

      prepared_data_name <- paste0("prepared_", data_name)
      if (!prepared_data_name %in% names(task)){
        log_warn("{data_name} is missing or has not been prepared yet")
        return(task)
      }

      log_info("Model is being applied on {prepared_data_name}")

      prediction <- predict_model(
        model = task[["model"]],
        new_data = task[[prepared_data_name]]
        )

      dup_names  <-  intersect(names(prediction %>% select(-siret, -periode)),
        names(task[[data_name]]))
      task[[data_name]]  <- task[[data_name]] %>% select(-one_of(dup_names))
      task[[data_name]] <- task[[data_name]] %>%
        left_join(prediction, by = c("siret", "periode"))

      log_info("Prediction successfully done.")
      return(task)
    }

    for (name in data_names){
      task  <- predict_on_given_data(name, task)
    }

    return(task)
}

#' Export de données
#'
#' Exporte les données.
#'
#' @inheritParams generic_task
#' @param export_type `"csv" | "mongodb"` \cr Export en csv, ou en
#'   mongodb ?
#' @param ... additional parameters for export functions.
#' @return `sf_task` \cr L'objet `task` donné en entrée.
#' @export
export.sf_task <- function(task, export_type, batch, ...){
  require(purrr)
  export_fields <- c(
    "siret",
    "periode",
    "raison_sociale",
    "departement",
    "region",
    "score",
    "score_diff",
    "connu",
    "date_ccsf",
    "etat_proc_collective",
    "date_proc_collective",
    "interessante_urssaf",
    # "default_urssaf",
    "effectif",
    "libelle_naf",
    "libelle_ape5",
    "code_ape",
    "montant_part_ouvriere",
    "montant_part_patronale",
    "ca",
    "ca_past_1",
    "benefice_ou_perte",
    "benefice_ou_perte_past_1",
    "resultat_expl",
    "resultat_expl_past_1",
    "poids_frng",
    "taux_marge",
    "frais_financier",
    "financier_court_terme",
    "delai_fournisseur",
    "dette_fiscale",
    "apart_heures_consommees",
    "apart_heures_autorisees",
    # "cotisation_moy12m",
    "compte_urssaf",
    "montant_majorations",
    "exercice_bdf",
    "exercice_diane",
    "delai"
    )

  f_scores <- c(F1 = 0.15, F2 = 0.07) # TODO TODO

  if (!is.null(export_type)) {
    assertthat::assert_that(all(export_type %in% c("csv", "mongodb")))

    log_info("Adding additional fields for export")

    res <- task[["new_data"]] %>%
      format_for_export(
        export_fields = export_fields,
        database = task[["database"]],
        collection = task[["collection"]],
        last_batch = batch
        )

    log_info("Data is exported to {paste(export_type, collapse = ' and ')}")
    purrr::walk(
      .x = export_type,
      .f = function(x, ...){
        if (x == "csv") {
          export_scores_to_csv(
            ...,
            relative_path = "../output/"
            )
        } else if (x == "mongodb") {
          export_scores_to_mongodb(
            ...,
            f_scores = f_scores,
            database = database,
            collection = "Scores"
            )
      }},
      formatted_data = res,
      batch = batch,
      algo = "algo"
      )
  }
  log_info("Data exported with success to
    {paste(export_type, collapse = ' and ')}")
    return(task)
}

#' Évaluation du modèle
#'
#' Evalue le modèle stocké dans l'objet task
#'
#' @inheritParams generic_task
#' @param eval_function `MLsegmentr::eval_function()` \cr Objet s3
#' d'évaluation.
#' @param data_name `character(1)` \cr Sur quelles données évaluer ?
#' @param plot `logical(1)` \cr Faut-il tracer la figure avec la fonction de
#'   l'`eval_function` (si disponible)
#' @param remove_strong_signals `logical(1)`\cr
#'   Faut-il retirer des échantillons de test ou de
#'   validation les entrerprises qui présentent des signaux forts, c'est-à-dire 3 mois de défaut, ou une
#'   procédure collective en cours ? Nécessite que les données contenues dans
#'   \code{task[["hist_data"]]} possèdent le champs "time_til_outcome".
#'
#' @describeIn evaluate
#'
#' @return `task` donnée en entrée à laquelle s'est ajoutée (ou a été
#' remplacé) un champs "model_performance", avec le résultat de la fonction
#' d'évaluation
#' @export
evaluate.sf_task <- function(
  task,
  eval_function =  NULL,
  data_name = c("validation_data"),
  plot = TRUE,
  prediction_names = "score",
  target_names = "outcome",
  segment_names = NULL,
  remove_strong_signals = TRUE,
  tracker = task[["tracker"]]
  ){
  if (is.null(eval_function)){
    default_eval_fun = TRUE
    eval_fun = MLsegmentr::eval_precision_recall()
  } else {
    default_eval_fun = FALSE
  }

  assertthat::assert_that(
    length(data_name) == 1,
    msg = "Evaluation can only be made on a single data.frame at once"
    )

  evaluation_data  <- task[[data_name]]
  if (remove_strong_signals){
    log_info("Les 'signaux forts' sont retires des donnees d'evaluation (test,
      validation)")
    # TODO add to log
    # task a7368366-ad3c-4dcb-b8d9-2e23de616bf0
    assertthat::assert_that("time_til_outcome" %in% names(evaluation_data))
    evaluation_data  <- evaluation_data %>%
      filter(is.na(time_til_outcome) | time_til_outcome > 0)
  }

  require(MLsegmentr)
  assesser <- Assesser$new(
    evaluation_data
    )

  assesser$set_predictions(prediction_names)
  assesser$set_targets(target_names)

  if (!is.null(segment_names)){
    assesser$set_segments(segment_names)
  }

  assesser$evaluation_funs <- eval_function

  perf <- assesser$assess_model(plot = plot)

  if (default_eval_fun){
    task[["model_performance"]] <- perf %>%
      filter(evaluation_name != "prcurve")
  } else {
    task[["model_performance"]] <- perf
  }

  if (!is.null(tracker)){
    tracker$set(
      model_performance = task[["model_performance"]]
    )
  }
  return(task)
  ## Log model performance.
}


#' Log des expérimentations Machine Learning
#'
#' Permet de loguer dans une collection mongodb la nature et les différents
#' éléments de la tâche d'apprentissage. Cf le package `MLlogs`.
#'
#' @inheritParams generic_task
#' @inheritParams mongodb_connection
#' @param collection `character(1)`. La collection vers laquelle exporter. Par
#' défaut, "ml_logs".
#'
#' @return returns the `task` unchanged.
#' @export
log.sf_task <- function(
  task,
  database = task[["database"]],
  collection = "ml_logs",
  tracker = task[["tracker"]],
  ...
  ){

  require("MLlogr")

  assertthat::assert_that(!is.null(tracker))

  tracker$database <- database
  tracker$collection <- collection
  tracker$set(
    model_parameters = task[["model_parameters"]],
    model_features =  task[["features"]],
    test_frame = task[["validation_data"]]
    )
  # No train data as it currently not work with cross validation.
  # e82f5df1-7c7f-472e-b3f8-60707d42136d

  tracker$log(...)

  return(invisible(task))
}

explain.sf_task <- function(task, ...){
  # Later
}

print.sf_task <- function(x, ...){
  require(purrr)
  cat("-- FIELDS --\n")
  aux_fun <- function(name, x){
    if (!is.character(x) || length(x) > 1){
      cat(paste0("  * ", name, " (", paste0(class(x), collapse = ", "), ")\n"))
    } else {
      cat(paste0("  * ", name, " : ", x, "\n"))
    }
  }
  purrr::walk2(names(x), x, aux_fun)

  cat("-- INFO --\n")
  purrr::walk2(
    names(x[["tracker"]]$values),
    x[["tracker"]]$values,
    aux_fun
    )
  return()
}
