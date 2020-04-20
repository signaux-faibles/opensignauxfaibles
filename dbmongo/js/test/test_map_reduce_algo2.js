

const notreMap = (testData) => {
  const results = {};
  emit = (key, value) => results[key] = value;
  map.call(testData); // will call emit an inderminate number of times
  // testData contains _id and value properties. testData is passed as this
  return results;
};

debug(JSON.stringify(notreMap(testData[0]), null, 2));
