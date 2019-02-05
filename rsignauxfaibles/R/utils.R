elapsed_months <- function(end_date, start_date) {
  ed <- as.POSIXlt(end_date)
  sd <- as.POSIXlt(start_date)
  12 * (ed$year - sd$year) + (ed$mon - sd$mon)
}


compare_vects <- function(data,field1,field2){
  
  smy <- data %>%
    group_by(siret) %>%
    summarize(a = first(!!sym(field1)), b = first(!!sym(field2)))
  return(list(only_first = smy$a & !smy$b, only_second = smy$b & !smy$a))
}


count_etab_entr <- function(df){
  nb_etab   <- n_distinct(df %>% select(siret))
  nb_entr <- n_distinct(df %>% select(siren))
  
  cat(nb_etab,' établissements dans ', nb_entr, ' entreprises','\n')
}


replace_na_by <- function(name,data,na_value) {
  data[is.na(data[,name]), name] = na_value
  return(data)
}

average_12m <- function(vec){
  sapply(
    1:length(vec),
    FUN = function(x){
      res <- tail(vec[1:x], 12)
      if (sum(!is.na(res)) > 3)
        return(mean(res , na.rm = TRUE))
      else return(NA)
    }
  )
}
