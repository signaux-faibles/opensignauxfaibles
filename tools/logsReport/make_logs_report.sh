#!/usr/bin/env Rscript
library("optparse")

option_list = list(
  make_option(c("-s", "--start_date"), type="character", default="2000-01-01T00:00:00Z", 
              help="Start date for logs in ISO-8601 format [default= %default].", metavar="character"),
  make_option(c("-p", "--parser"), type="character", default="all", 
              help="Which parser to select [default= %default]", metavar="character"),
  make_option(c("-i", "--interactive"), action = "store_true", type="logical", default = FALSE,
              help="Interactive mode flag", metavar="logical"),
  make_option(c("--include_fatal"), action = "store_true", type="logical", default = FALSE,
              help="Include fatal errors in report.", metavar="logical"),
  make_option(c("--debug"), action = "store_true", type="logical", default = FALSE,
              help="Run in debug mode (show code along with output)", metavar="logical")
); 

opt_parser = OptionParser(option_list=option_list);
opt = parse_args(opt_parser);

rmarkdown::render("make_report.Rmd",
                  encoding = "UTF-8",
                  params = list(
                    startdate = opt$start_date,
                    parsertype = opt$parser,
                    interactive = opt$interactive,
                    include_fatal = opt$include_fatal,
                    debug_mode = opt$debug
                  )
                )