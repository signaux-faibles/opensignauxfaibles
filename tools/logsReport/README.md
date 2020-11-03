# Integration Logs Report

## Set up

Start by making sure that you have R along with all needed poackages installed.

```
[vincent.viers@signfaib-labtenant opensignauxfaibles]$ R --version

R version 3.6.0 (2019-04-26) -- "Planting of a Tree"
Copyright (C) 2019 The R Foundation for Statistical Computing
Platform: x86_64-redhat-linux-gnu (64-bit)
```

Open the `R` REPL and type
```r
packages <- c(
  "tidyverse",
  "mongolite",
  "rjson",
  "DT",
  "kableExtra",
  "knitr",
  "rmarkdown"
)
packagecheck <- match(packages, utils::installed.packages()[,1])
packagestoinstall <- packages[is.na(packagecheck)]
install.packages(packagestoinstall)
```

## Generate report

The main R command to generate the report is :

```r
rmarkdown::render(
  "signauxfaibles/logs_reports/make_report.Rmd",
  encoding = "UTF-8",
  params = list(include_fatal = FALSE)
)
```

or, from the terminal :

```
R -e 'rmarkdown::render("signauxfaibles/logs_reports/make_report.Rmd", encoding = "UTF-8", params = list(startdate = "2020-10-28T00:00:00Z"))'
```

You can change more parameters by editing the `params` list argument. The available parameters along with their default values are:

- startdate: "2020-10-28T00:00:00Z"
- parsertype:  "all"
- include_fatal: FALSE
- interactive: FALSE
- mongo_collection: "Journal"
- mongo_db: "test"
- mongo_url: "mongodb://labbdd"
- debug_mode: FALSE




