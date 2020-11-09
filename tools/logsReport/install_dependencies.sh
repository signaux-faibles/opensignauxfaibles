#!/usr/bin/env Rscript
packages <- c(
  "tidyverse",
  "mongolite",
  "rjson",
  "DT",
  "kableExtra",
  "knitr",
  "rmarkdown",
  "optparse"
)
packagecheck <- match(packages, utils::installed.packages()[,1])
packagestoinstall <- packages[is.na(packagecheck)]
install.packages(packagestoinstall, repos = "https://cloud.r-project.org")

