package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
)

func NewSearchmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.SearchOpts,
		Short: "This command will allow you to perform search action on the resources provided based on actions provided.",
		Long:  `This command will allow you to perform search action on the resources provided based on actions provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cobra.CheckErr(fmt.Errorf("no resource provided for this command, please select resource from result of this command -> " + pkg.APPNAME + " " + explainCmd))
			}
		},
	}
	cmd.AddCommand(NewDatameshCmd(), NewDataSourceCmd())
	return cmd
}
