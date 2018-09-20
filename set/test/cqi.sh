echo 'INSERT INTO testbucket ( KEY, VALUE )
VALUES
(
 $id,
 {"id":$id}
)
' | cqi write -id=awesomeid
cqi read < query/selectall.n1ql
