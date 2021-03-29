# where_parser   
- percer in main.go
# Run the test query app
- go run main.go

# func in app
- Parse(query string, qb squirrel.SelectBuilder) (*squirrel.SelectBuilder, error)

simple parser for WHERE part of SQL queries

- Find(cond string) (string, []string)

search and return of the operator and parts of the condition

- splitCondition(condition string) squirrel.Sqlizer

replaces a single condition with a condition from the squirrel package

- splitAndCodition(andCond []string)  squirrel.Sqlizer

a common condition is created with "and"

- checkAnotherOpers(query string) (flag bool)

check for extra SQL statements in a query

- func callBackHandler(name string, value string) bool

checking whether the column type and the condition value match
