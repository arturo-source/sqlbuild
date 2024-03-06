# Generate SQL code only with your Golang structs

The main idea of this package is to generate the most standard SQL queries that exist. It is very easily understood with an example, the typical queries would be the following:

```go
package main

import (
 "fmt"

 "github.com/arturo-source/sqlbuild"
)

type Person struct {
 Id   int `db:"id"`
 Name string
 Age  *int `db:"age"`
}

func main() {
 age := 20
 p := Person{
  Id:   0,
  Name: "John",
  Age:  &age,
 }

 fmt.Println(sqlbuild.Create(p))
 fmt.Println(sqlbuild.SelectAll(p))
 fmt.Println(sqlbuild.SelectById(p))
 fmt.Println(sqlbuild.Insert(p))
 fmt.Println(sqlbuild.Update(p))
 fmt.Println(sqlbuild.DeleteAll(p))
 fmt.Println(sqlbuild.DeleteById(p))
 fmt.Println(sqlbuild.Drop(p))

 people := []Person{
  p,
  {Id: 1, Name: "Mike", Age: nil},
  {Id: 2, Name: "Cris", Age: nil},
 }
 fmt.Println(sqlbuild.InsertMultiple(people))
}
```

And the result of this code would be the following (I omit the err=nil):

```sql
CREATE TABLE "Person" ("id" INT NOT NULL PRIMARY KEY, "Name" TEXT NOT NULL, "age" INT)
SELECT * FROM "Person"
SELECT * FROM "Person" WHERE "id" = 0
INSERT INTO "Person" ("id", "Name", "age") VALUES (0, 'John', 20)
UPDATE "Person" SET "id" = 0, "Name" = 'John', "age" = 20 WHERE "id" = 0
DELETE FROM "Person"
DELETE FROM "Person" WHERE "id" = 0
DROP TABLE "Person"
INSERT INTO "Person" ("id", "Name", "age") VALUES (0, 'John', 20), (1, 'Mike', null), (2, 'Cris', null)
```

## TO DO

These are some problems I have to think how to face them:

- AUTO_INCREMENT is written differently in each database engine.
- References between tables. `CreateWithReferences` ??
- I don't really know how to avoid SQL injections. The package uses `sanitize` but I think it is not enough.

## Possible uses of this package

The reality is that I have created this package to practice using reflection, without a clear use case. But something came to my mind while I was programming it. If you can think others, don't hesitate to add it in an issue or PR:

1. Creation of SQL scripts using Go, for databases in development environments.
