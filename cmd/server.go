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
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("./templates/*")

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/gen", GenData)
	router.Run(":8080")
}

func GenData(c *gin.Context) {
	n := "./"
	builder := parser.NewParser()
	pkgs := builder.Parse(n)
	g := parser.NewGraph()
	c.JSON(http.StatusOK, g.Data(pkgs))
}
