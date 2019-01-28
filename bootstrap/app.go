package bootstrap

import (
	"extension/database"
	"fmt"
	"golang-laravel/app/exceptions"
	"golang-laravel/config"
	"golang-laravel/route"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-contrib/location"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func Start() {

	r := gin.Default()

	r.Static("/assets", "public")
	r.Use(exceptions.Message())
	r.Use(location.Default())

	/**********************************
	 * Initialize database connection *
	 **********************************/
	for _, databaseCfg := range config.Get().DATABASE {
		database.GetConnectionByDriver(databaseCfg.DRIVER).InitDB(map[string]config.Database{
			"default": databaseCfg,
		})
	}

	/*******************************************************************************************************
	 * Initialize template render and generate our templates map from our layouts/ and blocks/ directories *
	 *******************************************************************************************************/
	render := multitemplate.NewRenderer()
	templatesDir := config.Get().VIEW.PATH
	allBlocks, allLayouts := make([]string, 0), make([]string, 0)

	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if strings.HasPrefix(path, templatesDir+"/layouts/") {
				allLayouts = append(allLayouts, path)
			} else if strings.HasPrefix(path, templatesDir+"/blocks/") {
				allBlocks = append(allBlocks, path)
			}
		}

		return nil
	})

	if err != nil {
		panic(err.Error())
	}

	for _, block := range allBlocks {
		viewCopy := make([]string, 0)
		viewCopy = append(viewCopy, block)

		for _, layout := range allLayouts {
			viewCopy = append(viewCopy, layout)
		}

		identify := strings.Replace(block, config.Get().VIEW.PATH+"/blocks/", "", -1)
		identify = strings.Replace(identify, ".html", "", -1)
		identify = strings.Replace(identify, "/", ".", -1)

		render.AddFromFiles(identify, viewCopy...)
	}

	r.HTMLRender = render

	/******************************
	 * Initialize custom uri path *
	 ******************************/
	route.Setup(r)

	r.Run(":" + strconv.Itoa(config.Get().SERVER.PORT))
}
