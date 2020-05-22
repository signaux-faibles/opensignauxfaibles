#!/bin/bash

if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' 's/^export //' ./**/*.js
else
  sed -i 's/^export //' ./**/*.js
fi
