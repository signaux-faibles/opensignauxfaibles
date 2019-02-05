quantile_APE <- function(data, variable_names, levels = 1) {

  assertthat::assert_that(all(c('code_naf_niveau1', 'code_ape') %in% names(data)))
  assertthat::assert_that(all(levels >= 1 & levels <= 5))

  # FIX ME: problème de nom dans le mutate_at si une seule variable
  assertthat::assert_that(length(variable_names) > 1)

  for (i in seq_along(levels)) {
    level = levels[i]

    if (level == 1) {
      data <- data %>%
        mutate(.target = as.factor(code_naf_niveau1))
    } else {
      data <- data %>%
        mutate(.target = as.factor(substr(code_ape, 1, levels)))
    }

    aux_fun <- function(x){
      if (all(is.na(x))) return (NA_real_)
      mecdf = ecdf(x)(x)
      return(mecdf)
    }

    my_fun <- c("aux_fun")
    names(my_fun) <- paste0("distrib_APE",level)

    data <- data %>%
      group_by(.target) %>%
      mutate_at(
        .vars = variable_names,
        .funs = funs_(my_fun)) %>%
      ungroup() %>%
      select(-.target)
  }

  return(data)


}

quantile_APE_create <- function(
  my_data,
  variable_names,
  ape_levels = 1
){
  
  assertthat::assert_that(all(c("code_naf", "code_ape") %in%
                                names(my_data)))
  
  assertthat::assert_that(all(ape_levels >= 1 & ape_levels <= 5))
  
  
  out <- data.frame(
    niveau = numeric(0),
    code = character(0),
    variable = character(0),
    moy = numeric(0),
    std = numeric(0)
  )
  
  for (i in seq_along(ape_levels)) {
    level <- ape_levels[i]
    
    if (level == 1) {
      my_data <- my_data %>%
        mutate(.target = as.character(code_naf))
    } else {
      my_data <- my_data %>%
        mutate(.target = as.character(substr(code_ape, 1, ape_levels)))
    }
    
    for (i in seq_along(variable_names)){
      out_add  <- my_data %>%
        group_by(.target) %>%
        summarize(
          niveau = level,
          code = unique(.target),
          variable = variable_names[i],
          moy = mean(!!rlang::sym(variable_names[i]), na.rm = TRUE),
          std = sd(!!rlang::sym(variable_names[i]), na.rm = TRUE)
        )
      out <- rbind(out, out_add)
    }
  }
  
  out <- out %>%
    select(-c(".target"))
  
  out[!is.finite(out$moy), 'moy'] <- NA
  out[!is.finite(out$std), 'std'] <- NA
  return(out)
}

quantile_APE_apply <- function(
  ref_quantile_APE, ...){
  
  assertthat::assert_that(all(
    c("niveau",
      "code",
      "variable",
      "moy",
      "std") %in%
      names(ref_quantile_APE)
  ))
  
  ape_levels  <- unique(ref_quantile_APE$niveau)
  variable_names <- unique(ref_quantile_APE$variable)
  
  aux_fun  <- function(my_data){
    
    out  <- my_data
    
    for (i in seq_along(ape_levels)) {
      level <- ape_levels[i]
      
      if (level == 1) {
        out <- out %>%
          mutate(.target = as.character(code_naf))
      } else {
        out <- out %>%
          mutate(.target = as.character(substr(code_ape, 1, ape_levels)))
      }
      
      
      for (i in seq_along(variable_names)){
        
        new_var  <- paste0(variable_names[i], "_distrib_APE", level)
        out  <- out %>%
          left_join(ref_quantile_APE %>%
                      filter(variable == variable_names[i]) %>%
                      select(code, moy, std),
                    by = c(".target" = "code")) %>%
          mutate(!!new_var := (!!rlang::sym(variable_names[i]) - moy) / std) %>%
          select(-c("moy", "std"))
      }
      out <- out %>%
        select(-.target)
    }
    return(out)
  }
  
  return(
    lapply(..., aux_fun)
  )
}


