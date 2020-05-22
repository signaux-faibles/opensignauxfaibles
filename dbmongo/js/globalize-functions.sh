#!/bin/bash

# We use perl because sed adds an empty line at the end of every js file,
# which was adding changes to git's staging, while debugging failing tests.
perl -pi'' -e 's/^export //' ./**/*.js
