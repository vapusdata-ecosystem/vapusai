package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/platform/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai-studio/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
)

func NewUserCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   userResource,
		Short: "This command is interface to interact with the platform for datamesh resources.",
		Long:  `This command is interface to interact with the platform for datamesh resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			vapusGlobals.VapusPlatformClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(currentAccessToken),
				Action:      userActions(cmd.Parent().Use),
			}
			err := vapusGlobals.VapusPlatformClient.HandleAction()
			if err != nil {
				cobra.CheckErr(err)
			}

			defer vapusGlobals.VapusPlatformClient.Close()
		},
	}
	return cmd
}

func userActions(parentCmd string) string {
	switch parentCmd {
	case pkg.GetOps:
		return pb.UserAgentOperations_GET_USER.String()
	case pkg.DescribeOps:
		return pb.UserAgentOperations_GET_USER.String()
	default:
		return ""
	}
}

// func getuser() {
// 	err := vapusGlobals.VapusPlatformClient.ListActions(pb.UserAgentOperations_GET_USER.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeUser() {

// 	err := vapusGlobals.VapusPlatformClient.DescribeActions(pb.UserAgentOperations_GET_USER.String(), viper.GetString(currentAccessToken), "")
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
