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

## Main problems

There are some problems that I have not been able to standardize, such as **AUTO_INCREMENT is written differently in each database engine**, and I wanted this package to be database agnostic. So the CREATE TABLE does not contain a clause to increment the id.

Another problem that I'm thinking about how to solve is **references between tables**, because obviously, it is very common in SQL databases to have relationships between tables.

And the last one is that **I don't really know how to avoid SQL injections**, so I created the `sanitize` function, but this only duplicates quotes.

## Possible uses of this package

The reality is that I have created this package to practice using reflection, without a clear use case. But something came to my mind while I was programming it. If you can think others, don't hesitate to add it in an issue or PR:

1. Creation of SQL scripts using Go, for databases in development environments.
