package route

import "extension/reflection"

func Setup(route reflection.WebRoute) {

	// e.Use(middleware.IsAuth)

	e.GET("/", "Home@Display")
	e.GET("/home", "Home@Display")
}
