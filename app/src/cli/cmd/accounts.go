package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/platform/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai-studio/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
)

func NewAccountCmd() *cobra.Command {
	var action, file string
	var la bool
	cmd := &cobra.Command{
		Use:   accountResource,
		Short: "This command is interface to interact with the platform for datamesh resources.",
		Long:  `This command is interface to interact with the platform for datamesh resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			if la {
				vapusGlobals.VapusPlatformClient.ListResourceActions("account")
				return
			}
			resAct := getAccountAction(cmd.Parent().Use, action)
			spinner := pkg.GetSpinner(36)
			spinner.Start()
			vapusGlobals.VapusPlatformClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(currentAccessToken),
				Action:      resAct,
				File:        file,
			}
			err := vapusGlobals.VapusPlatformClient.HandleAction()
			spinner.Stop()
			if err != nil {
				cobra.CheckErr(err)
			}

			defer vapusGlobals.VapusPlatformClient.Close()

		},
	}
	cmd.PersistentFlags().StringVar(&action, "action", "", "Action for the platform that should be executed on current resource with params in a file")
	cmd.PersistentFlags().StringVar(&file, "file", "", "File containing the parameters for the action")
	cmd.PersistentFlags().BoolVar(&la, "la", false, "List down all the actions that can be performed on the current resource")

	return cmd
}

func getAccountAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.GetOps:
		return pb.AccountAgentActions_LIST_ACCOUNT.String()
	case pkg.DescribeOps:
		return pb.AccountAgentActions_LIST_ACCOUNT.String()
	case pkg.ActOps:
		return action
	default:
		return pkg.ErrInvalidAction.Error()
	}
}

// func accountActions(parentCmd string) {
// 	switch parentCmd {
// 	case pkg.GetOps:
// 		getAccount()
// 	case pkg.DescribeOps:
// 		describeAccount()
// 	default:
// 		cobra.CheckErr("Invalid action")
// 	}
// }

// func getAccount() {
// 	err := vapusGlobals.VapusPlatformClient.ListActions(pb.AccountAgentActions_LIST_ACCOUNT.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeAccount() {
// 	err := vapusGlobals.VapusPlatformClient.DescribeActions(pb.AccountAgentActions_LIST_ACCOUNT.String(), viper.GetString(currentAccessToken), "")
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
