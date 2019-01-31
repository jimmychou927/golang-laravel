# Golang-Laravel

## Contents

- [Intro](#intro)
- [Installation](#installation)
- [Route Examples](#route-examples)
- [ORM Examples](#orm-examples)
- [Exception](#exception)
- [Message](#message)

## Intro
Based on Gin Web Framework and Go-Admin modules
Gin Web Framework https://github.com/gin-gonic/gin
GoAdmin https://github.com/chenhg5/go-admin

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
Based on Go-Admin (https://github.com/chenhg5/go-admin) db module
### Find
```go
item, _ := ModelName.Orm().Find(1)
fmt.Println(item['id'])
fmt.Println(item['name'])
```

### Where
```go
// with operation
result, _ := ModelName.Orm().Where("id", "=", 1).All()

// without operation
result, _ := ModelName.Orm().Where("id", 1).All()

// fetch result
for idx, value := range result {
    // do something ...
}
```

### WhereIn/WhereNotIn
```go
// with integer array
idList := []int{1, 2, 3, 4, 5, 6}
result, _ := ModelName.Orm().WhereIn("id", idList).All()

// with string array
idList := []string{"a", "b", "c"}
result, _ := ModelName.Orm().WhereIn("id", idList).All()

// with interface array
idList := []interface{}{1, "a", "b"}
result, _ := ModelName.Orm().WhereIn("id", idList).All()

// with sub query
subQuery := SubModelName.Orm().Where("id", 3)
result, _ := MainModelName.Orm().WhereIn("master_id", subQuery).All()

// fetch result
for idx, value := range result {
    // do something ...
}
```

### WhereRaw
```go
result1, _ := ModelName.Orm().WhereRaw("id = ?", 21).All()
result2, _ := ModelName.Orm().WhereRaw("id IN (?, ?)", 21, 43).All()

// fetch result
for idx, value := range result {
    // do something ...
}
```

### First
```go
item, _ := ModelName.Orm().Where("name", "golang").First()
fmt.Println(item['id'])
fmt.Println(item['name'])
```

### Count
```go
total, _ := ModelName.Orm().Where("id", 1).Count()
fmt.Println(total)
```

### Select
```go
_, _ = ModelName.Orm().Select("id", "first_name", "last_name", "pwd").All()
```

### SelectRaw
```go
_, _ = ModelName.Orm().SelectRaw("pwd as password").SelectRaw("concat(first_name, ' ', last_name) as name").All()
```

### Join/LeftJoin/RightJoin
```go
// without table alias
result, _ := ModelName.Orm().Join("sub_table", "sub_table.main_id", "=", "main_table.id").All()

// with table alias
result, _ := ModelName.Orm().Join("sub_table", "sub_table_alias", "sub_table_alias.main_id", "=", "main_table.id").All()

// with sub query
subQuery := SubModelName.Orm().Where("type", 3)
result, _ := MainModelName.Orm().Join(subQuery, "sub_alias", "sub_alias.main_id", "=", "main_model.id")

// fetch result
for idx, value := range result {
    // do something ...
}
```

### GroupBy
```go
_, _ = ModelName.Orm().Where("type", 3).GroupBy("status")
_, _ = ModelName.Orm().Where("type", 3).GroupBy("status", "user_id", "date")

// fetch result
for idx, value := range result {
    // do something ...
}
```

### OrderBy
```go
_, _ = ModelName.Orm().Where("type", 3).OrderBy("status", "asc").OrderBy("user_id", "desc").OrderBy("date", "desc")

// fetch result
for idx, value := range result {
    // do something ...
}
```

### Exception
Throw console message to browser

### Message
Display variable content on browser
```go
import "extension/print"

func main() {
    data := map[string]map[string]interface{}{
        "header": {
            "height": 200,
            "content": "my title",
        },
        "content": {
            "height": 500,
            "content": "my content",
        },
        "footer": {
            "height": 300,
            "content": "my footer",
        },
    }

    print.Msg(data)
}
```
