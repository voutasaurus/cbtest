set -ex
curl http://localhost:8080/new -d '{"id":"1","doc":{"type":"x","x":"123lol"}}'
curl http://localhost:8080/new -d '{"id":"2","doc":{"type":"x","x":"123kek"}}'
curl http://localhost:8080/search -d '{"doc":{"x":""}}'
set +ex
