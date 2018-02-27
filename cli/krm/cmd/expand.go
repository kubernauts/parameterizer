package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/kubernauts/parameterizer/pkg/parameterizer"
	"github.com/spf13/cobra"
)

// expandCmd represents the expand command
var expandCmd = &cobra.Command{
	Use:   "expand",
	Short: "Expand an app definition to a YAML manifest",
	Long: `Takes an Parameterizer YAML manifest and creates a Kubernetes YAML
manifest that you can feed into an installer.

For example:

$ krm expand install-ghost-with-helm.yaml | kubectl apply -f -`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := parameterizer.Parse(args[0])
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Println(p)
	},
}

func init() {
	rootCmd.AddCommand(expandCmd)
}
