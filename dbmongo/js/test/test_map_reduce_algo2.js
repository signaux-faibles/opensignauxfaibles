

const notreMap = (testData) => {
  const results = [];
  emit = (result) => results.push(result);
  map(testData); // will call emit an inderminate number of times
  return results;
};

debug(JSON.stringify(notreMap(testData), null, 2));
