package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/platform/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai-studio/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
)

func NewDatameshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   datameshResource,
		Short: "This command is interface to interact with the platform for datamesh resources.",
		Long:  `This command is interface to interact with the platform for datamesh resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			if la {
				vapusGlobals.VapusPlatformClient.ListResourceActions("datamesh")
				return
			}
			resAct := getDatameshAction(cmd.Parent().Use, action)
			spinner := pkg.GetSpinner(36)
			spinner.FinalMSG = "Done"
			spinner.Start()
			vapusGlobals.VapusPlatformClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(currentAccessToken),
				Action:      resAct,
				File:        file,
				SearchQ:     search,
			}
			err := vapusGlobals.VapusPlatformClient.HandleAction()
			spinner.Stop()
			if err != nil {
				cobra.CheckErr(err)
			}

			defer vapusGlobals.VapusPlatformClient.Close()
		},
	}
	return cmd
}

func getDatameshAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.GetOps:
		return pb.DataMeshAgentActions_LIST_DATAMESH.String()
	case pkg.DescribeOps:
		return pb.DataMeshAgentActions_LIST_DATAMESH.String()
	case pkg.ActOps:
		return action
	case pkg.SearchOpts:
		return pb.VapusSearchType_DATAPRODUCTS.String()
	default:
		return pkg.ErrInvalidAction.Error()
	}
}

// func datameshActions(parentCmd string) {
// 	switch parentCmd {
// 	case pkg.GetOps:
// 		getDatamesh()
// 	case pkg.DescribeOps:
// 		describeDatamesh()
// 	default:
// 		cobra.CheckErr("Invalid action")
// 	}
// }

// func getDatamesh() {
// 	err := vapusGlobals.VapusPlatformClient.ListActions(pb.DataMeshAgentActions_LIST_DATAMESH.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeDatamesh() {
// 	err := vapusGlobals.VapusPlatformClient.DescribeActions(pb.DataMeshAgentActions_LIST_DATAMESH.String(), viper.GetString(currentAccessToken), "")
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
