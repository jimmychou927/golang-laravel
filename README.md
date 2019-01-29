# Golang-Laravel

## Contents

- [Installation](#installation)
- [ORM Examples](#orm-examples)

## Installation
```
cd $GOPATH/src
git clone https://github.com/jimmychou927/golang-laravel
```

## ORM Examples

### Find
```go
item, _ := ModelName.Orm().Find(1)
fmt.Println(item['id'])
fmt.Println(item['name'])
```

### Where
```go
result, _ := ModelName.Orm().Where("id", "=", 1).All()
for idx, value := range result {
    // do something ...
}
```

### WhereIn/WhereNotIn
```go
idList := []interface{}{1, 2, 3, 4, 5, 6}
result, _ := ModelName.Orm().WhereIn("id", idList).All()
for idx, value := range result {
    // do something ...
}
```

### WhereRaw
```go
result1, _ := ModelName.Orm().WhereRaw("id = ?", 21).All()
result2, _ := ModelName.Orm().WhereRaw("id IN (?, ?)", 21, 43).All()
for idx, value := range result {
    // do something ...
}
```

### First
```go
item, _ := ModelName.Orm().Where("name", "=", "golang").First()
fmt.Println(item['id'])
fmt.Println(item['name'])
```

### Count
```go
total, _ := ModelName.Orm().Where("id", "=", 1).Count()
fmt.Println(total)
```

### Join/LeftJoin
```go
result, _ := ModelName.Orm().Join("sub_table", "sub_table.main_id", "=", "main_table.id").All()
for idx, value := range result {
    // do something ...
}
```

### WhereInQuery
```go
subQuery := SubModelName.Orm().Where("id", "=", 3)
result, _ := MainModelName.Orm().WhereInQuery("master_id", subQuery).All()
for idx, value := range result {
    // do something ...
}
```
