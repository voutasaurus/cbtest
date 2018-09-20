echo "== writing idA,idB,idC:"
curl http://localhost:8080/write -d '{"Stmt":"INSERT INTO testbucket ( KEY, VALUE ) VALUES ( $id, {\"id\":$id} )", "Args":{"id":"idA"}}'
curl http://localhost:8080/write -d '{"Stmt":"INSERT INTO testbucket ( KEY, VALUE ) VALUES ( $id, {\"id\":$id} )", "Args":{"id":"idB"}}'
curl http://localhost:8080/write -d '{"Stmt":"INSERT INTO testbucket ( KEY, VALUE ) VALUES ( $id, {\"id\":$id} )", "Args":{"id":"idC"}}'
echo "== without where:"
curl http://localhost:8080/read -d '{"Stmt":"SELECT *,META(t).id FROM testbucket as t"}'
echo "== with where:"
curl http://localhost:8080/read -d '{"Stmt":"SELECT *,META(t).id FROM testbucket as t WHERE id == $id", "Args":{"id":"idA"}}'
