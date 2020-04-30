# This file is run by dbmongo/js_test.go.

set -e

jsc ../public/*.js ../common/*.js public/lib_public.js objects.js public/test_public.js
