package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/platform/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai-studio/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
)

func NewDataCatalogCmd() *cobra.Command {
	var action, file string
	var la bool
	cmd := &cobra.Command{
		Use:   dataCatalogResource,
		Short: "This command is interface to interact with the platform for dataCatalog resources.",
		Long:  `This command is interface to interact with the platform for dataCatalog resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			if la {
				vapusGlobals.VapusPlatformClient.ListResourceActions("datacatalog")
				return
			}
			resAct := getDatacatalogAction(cmd.Parent().Use, action)
			spinner := pkg.GetSpinner(36)
			spinner.Prefix = "Performing " + resAct + " action for the current platform"
			spinner.Start()
			vapusGlobals.VapusPlatformClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(currentAccessToken),
				Action:      resAct,
				File:        file,
			}
			err := vapusGlobals.VapusPlatformClient.HandleAction()
			if err != nil {
				cobra.CheckErr(err)
			}

			defer vapusGlobals.VapusPlatformClient.Close()
			spinner.Stop()
		},
	}
	cmd.PersistentFlags().StringVar(&action, "action", "", "Action for the platform that should be executed on current resource with params in a file")
	cmd.PersistentFlags().StringVar(&file, "file", "", "File containing the parameters for the action")
	cmd.PersistentFlags().BoolVar(&la, "la", false, "List down all the actions that can be performed on the current resource")

	return cmd
}

func getDatacatalogAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.GetOps:
		return pb.DataCatalogAgentActions_LIST_CATALOG.String()
	case pkg.DescribeOps:
		return pb.DataCatalogAgentActions_GET_SELF_DATA_CATALOG.String()
	case pkg.ActOps:
		return action
	default:
		return pkg.ErrInvalidAction.Error()
	}
}

// func dataCatalogActions(parentCmd string) {
// 	switch parentCmd {
// 	case pkg.GetOps:
// 		getDataCatalog()
// 	case pkg.DescribeOps:
// 		describeDataCatalog()
// 	default:
// 		cobra.CheckErr("Invalid action")
// 	}
// }

// func getDataCatalog() {
// 	err := vapusGlobals.VapusPlatformClient.ListActions(pb.DataCatalogAgentActions_LIST_CATALOG.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeDataCatalog() {
// 	err := vapusGlobals.VapusPlatformClient.DescribeActions(pb.DataCatalogAgentActions_GET_SELF_DATA_CATALOG.String(), viper.GetString(currentAccessToken), "")
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
