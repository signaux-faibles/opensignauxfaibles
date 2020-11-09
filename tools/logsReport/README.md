# Integration Logs Report

## Set up

Start by making sure that you have R along with all needed packages installed. You will also need to install [`pandoc`](https://pandoc.org/installing.html) if you don't have it. Linux users may need to install other dependencies as well as compiler-related libraries such as `build-essential` or `gcc`.

```sh
$ R --version

R version 3.6.0 (2019-04-26) -- "Planting of a Tree"
Copyright (C) 2019 The R Foundation for Statistical Computing
Platform: x86_64-redhat-linux-gnu (64-bit)
```

Install R packages 

```sh
./install_dependencies.sh
```

## Generate report

The main R command to generate the report is :

```python
rmarkdown::render(
  "signauxfaibles/logs_reports/make_report.Rmd",
  encoding = "UTF-8",
  params = list(...)
)
```

Which was packaged in a convenient CLI `make_logs_report.sh`

```sh
$ ./make_logs_report.sh --help
Usage: ./make_logs_report.sh [options]


Options:
        -s CHARACTER, --start_date=CHARACTER
                Start date for logs in ISO-8601 format [default= 2000-01-01T00:00:00Z].

        -p CHARACTER, --parser=CHARACTER
                Which parser to select [default= all]

        -i, --interactive
                Interactive mode flag

        --include_fatal
                Include fatal errors in report.

        --debug
                Run in debug mode (show code along with output)

        -h, --help
                Show this help message and exit
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