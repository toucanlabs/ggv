package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/toucan-labs/ggv/internal/parser"
)

var versionCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func serve() {
	router := gin.Default()
	router.Static("/assets", "../assets")
	router.LoadHTMLGlob("./templates/*")

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/gen", func(c *gin.Context) {
		n := "./"
		builder := parser.NewParser()
		pkgs := builder.Parse(n)

		jsobResult := map[string]interface{}{}
		for _, p := range pkgs {
			jsobResult[p.Dir] = p.Funcs()
		}

		c.JSON(http.StatusOK, jsobResult)
	})
	router.Run(":8080")
}
