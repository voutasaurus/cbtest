set -ex
curl http://localhost:8080/new -d '{"key":"1","value":{"type":"x","x":"123lol"}}'
curl http://localhost:8080/new -d '{"key":"2","value":{"type":"x","x":"123kek"}}'
curl http://localhost:8080/search -d '{"value":{"x":""}}'
set +ex
