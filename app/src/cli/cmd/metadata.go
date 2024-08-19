package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	dpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/domain/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai-studio/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
)

func NewMetadataCmd() *cobra.Command {
	var dsId string
	cmd := &cobra.Command{
		Use:   pkg.MetaDataResource,
		Short: "This command is interface to interact with the platform for metadata resources.",
		Long:  `This command is interface to interact with the platform for metadata resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			defer vapusGlobals.VapusPlatformClient.Close()
			if la {
				vapusGlobals.VapusPlatformClient.ListResourceActions("metadata")
				return
			}
			resAct := getMetadataAction(cmd.Parent().Use, action)
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
			if dsId != "" {
				vapusGlobals.VapusPlatformClient.ActionHandler.Identifier = []string{dsId}
			}
			err := vapusGlobals.VapusPlatformClient.HandleAction()
			if err != nil {
				spinner.Stop()
				cobra.CheckErr(err)
			}

			spinner.Stop()

		},
	}
	cmd.PersistentFlags().StringVar(&dsId, "id", "", "Data product Id to perform the action on")
	return cmd
}

func getMetadataAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.GetOps:
		return dpb.DataSourceAgentActions_LIST_DATASOURCE.String()
	case pkg.DescribeOps:
		return dpb.DataSourceAgentActions_DESCRIBE_DATASOURCE.String()
	case pkg.ActOps:
		return action
	default:
		return pkg.ErrInvalidAction.Error()
	}
}

// func getDataSource() {
// 	err := vapusGlobals.VapusPlatformClient.ListActions(pb.DataSourceAgentActions_LIST_DATASOURCE.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeDataSource(args []string) {
// 	if len(args) < 1 {
// 		cobra.CheckErr("Invalid number of arguments, please provide the dataSource ID")
// 	}
// 	err := vapusGlobals.VapusPlatformClient.DescribeActions(pb.DataSourceAgentActions_DESCRIBE_DATASOURCE.String(), viper.GetString(currentAccessToken), args[0])
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
