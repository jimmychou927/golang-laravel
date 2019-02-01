package route

import "extension/reflection"

func Setup(route reflection.WebRoute) {

	// route.Use(middleware.IsAuth)

	route.GET("/", "Home@Display")
	route.GET("/home", "Home@Display")
}
