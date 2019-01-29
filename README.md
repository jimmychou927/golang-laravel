# Golang-Laravel

## Contents

- [Installation](#installation)
- [Route Examples](#route-examples)
- [ORM Examples](#orm-examples)

## Installation
```
cd $GOPATH/src
git clone https://github.com/jimmychou927/golang-laravel
```

## Route Examples
Router definition in file route/web.go
```go
route.GET("/", "Home@Display")

route.GET("/login", "Auth@Display")
route.POST("/", "Auth@Login")
```
Controller Auth in folder app/http/controller/auth.go
```go
package controller

import (
	"github.com/gin-gonic/gin"
)

type Auth struct {
	Controller
}

func (this *Auth) Display(c *gin.Context) {
    // do GET Request procedure ...
}

func (this *Auth) Login(c *gin.Context) {
    // do POST Request procedure ...
}

func init() {

	register(Auth{})
}
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

### Select
```go
```

### SelectRaw
```go
```

### Join/LeftJoin
```go
result, _ := ModelName.Orm().Join("sub_table", "sub_table.main_id", "=", "main_table.id").All()
for idx, value := range result {
    // do something ...
}
```
### JoinQuery/LeftJoinQuery
```go
```

### WhereInQuery
```go
subQuery := SubModelName.Orm().Where("id", "=", 3)
result, _ := MainModelName.Orm().WhereInQuery("master_id", subQuery).All()
for idx, value := range result {
    // do something ...
}
```
### GroupBy
```go
```

### OrderBy
```go
```
