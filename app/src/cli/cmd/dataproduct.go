package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	dpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/domain/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai-studio/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
)

func NewDataProductCmd() *cobra.Command {
	var dpId string
	cmd := &cobra.Command{
		Use:   dataProductResource,
		Short: "This command is interface to interact with the platform for dataProduct resources.",
		Long:  `This command is interface to interact with the platform for dataProduct resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			if la {
				vapusGlobals.VapusPlatformClient.ListResourceActions("dataproduct")
				return
			}
			resAct := getDataproductAction(cmd.Parent().Use, action)
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
			if dpId != "" {
				vapusGlobals.VapusPlatformClient.ActionHandler.Identifier = []string{dpId}
			}
			err := vapusGlobals.VapusPlatformClient.HandleAction()
			if err != nil {
				spinner.Stop()
				cobra.CheckErr(err)
			}

			defer vapusGlobals.VapusPlatformClient.Close()
			spinner.Stop()

		},
	}
	cmd.PersistentFlags().StringVar(&dpId, "id", "", "Data product Id to perform the action on")
	return cmd
}

func getDataproductAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.GetOps:
		return dpb.DataProductAgentActions_LIST_DATAPRODUCT.String()
	case pkg.DescribeOps:
		return dpb.DataProductAgentActions_DESCRIBE_DATAPRODUCT.String()
	case pkg.ActOps:
		return action
	case pkg.SearchOpts:
		return action
	default:
		return pkg.ErrInvalidAction.Error()
	}
}

// func dataProductActions(parentCmd string, args []string) {
// 	switch parentCmd {
// 	case pkg.GetOps:
// 		getDataProduct()
// 	case pkg.DescribeOps:
// 		describeDataProduct(args)
// 	default:
// 		cobra.CheckErr("Invalid action")
// 	}
// }

// func getDataProduct() {
// 	err := vapusGlobals.VapusPlatformClient.ListActions(pb.DataProductAgentActions_LIST_DATAPRODUCT.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeDataProduct(args []string) {
// 	if len(args) < 1 {
// 		cobra.CheckErr("Invalid number of arguments, please provide the dataProduct ID")
// 	}
// 	err := vapusGlobals.VapusPlatformClient.DescribeActions(pb.DataProductAgentActions_DESCRIBE_DATAPRODUCT.String(), viper.GetString(currentAccessToken), args[0])
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
