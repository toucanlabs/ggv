package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/toucan-labs/ggv/internal/parser"
)

var serverCommand = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

var pkgCfg string
var internalFunc bool

func init() {
	rootCmd.AddCommand(serverCommand)
	rootCmd.PersistentFlags().StringVar(&pkgCfg, "pkg", "./", "")
	rootCmd.PersistentFlags().BoolVar(&internalFunc, "internal", false, "")

}

func serve() {
	router := gin.Default()
	router.GET("/index", func(c *gin.Context) {
		contents := `
		<head>
  <style> body { margin: 0; } </style>

  <script src="https://bundle.run/@yarnpkg/lockfile@1.1.0"></script>
  <script src="https://unpkg.com/dat.gui"></script>
  <script src="https://unpkg.com/d3-quadtree"></script>
  <script src="https://unpkg.com/d3-force"></script>
  <script src="https://unpkg.com/force-graph@1.43.4/dist/force-graph.min.js"></script>
</head>

<body>
  <div id="graph"></div>

  <script>
    // controls
    const controls = { 'DAG Orientation': 'lr'};
    const gui = new dat.GUI();
    gui.add(controls, 'DAG Orientation', ['lr', 'td', 'radialout', null])
      .onChange(orientation => graph && graph.dagMode(orientation));

    // graph config
    const graph = ForceGraph()
      .backgroundColor('#101020')
      .linkColor(() => 'rgba(255,255,255,0.2)')
      .dagMode('lr')
      .dagLevelDistance(300)
      .nodeId('id')
      .linkCurvature(d =>
        0.07 * // max curvature
        // curve outwards from source, using gradual straightening within a margin of a few px
        (['td', 'bu'].includes(graph.dagMode())
          ? Math.max(-1, Math.min(1, (d.source.x - d.target.x) / 25)) :
          ['lr', 'rl'].includes(graph.dagMode())
            ? Math.max(-1, Math.min(1, (d.target.y - d.source.y) / 25))
            : ['radialout', 'radialin'].includes(graph.dagMode()) ? 0 : 1
        )
      )
      .linkDirectionalParticles(2)
      .linkDirectionalParticleWidth(3)
      .nodeCanvasObject((node, ctx) => {
        const label = node.val;
        const fontSize = 15;
        ctx.font = "15px Sans-Serif";
        const textWidth = ctx.measureText(label).width;
        const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.2); // some padding

        ctx.fillStyle = 'rgba(0, 0, 0, 0.2)';
        ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, ...bckgDimensions);

        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillStyle = 'lightsteelblue';
        ctx.fillText(label, node.x, node.y);

        node.__bckgDimensions = bckgDimensions; // to re-use in nodePointerAreaPaint
      })
      .nodePointerAreaPaint((node, color, ctx) => {
        ctx.fillStyle = color;
        const bckgDimensions = node.__bckgDimensions;
        bckgDimensions && ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, ...bckgDimensions);
      })
      .d3Force('collide', d3.forceCollide(13))
      .d3AlphaDecay(0.02)
      .d3VelocityDecay(0.3);

    fetch('/gen?pkg=name')
      .then(res => res.json()).then(data => {
          graph(document.getElementById('graph'))
          .graphData(data);
      });
  </script>
</body>
		`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(contents))
	})

	router.GET("/gen", GenData)
	router.Run(":8080")
}

func GenData(c *gin.Context) {
	n := pkgCfg
	builder := parser.NewParser()
	pkgs := builder.Parse(n)
	g := parser.NewGraph()
	c.JSON(http.StatusOK, g.Data(internalFunc, pkgs))
}
