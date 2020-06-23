(cd dbmongo/lib/engine && go generate .) && \
(cd dbmongo/js && npm run lint && npm test) && \
(cd dbmongo && go test) && \
(killall dbmongo; cd dbmongo && go build) && \
./test-api.sh && \
./test-api-2.sh
