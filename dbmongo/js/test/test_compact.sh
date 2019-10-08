result_add=$(jsc ../compact/currentState.js compact/test_current_state.js)
if [ "$result_add" != 'true' ]; then
  echo "$result_add"
  exit 1
fi

result_add=$(jsc ../compact/currentState.js ../compact/reduce.js ./testing.js ./compact/test_reduce.js)
if [ "$result_add" != 'true' ]; then
  echo "$result_add"
  exit 1
fi

