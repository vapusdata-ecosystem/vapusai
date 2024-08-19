package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	dpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/domain/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai-studio/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
)

func NewDataWorkerCmd() *cobra.Command {
	var action, file string
	var la bool
	cmd := &cobra.Command{
		Use:   dataWorkerResource,
		Short: "This command is interface to interact with the platform for dataWorker resources.",
		Long:  `This command is interface to interact with the platform for dataWorker resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			defer vapusGlobals.VapusPlatformClient.Close()
			if la {
				vapusGlobals.VapusPlatformClient.ListResourceActions("dataworker")
				return
			}
			resAct := getDataworkerAction(cmd.Parent().Use, action)
			spinner := pkg.GetSpinner(36)
			spinner.Prefix = "Performing " + resAct + " action for the current platform"
			spinner.Suffix = "Please wait...\n"
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
				spinner.Stop()
				cobra.CheckErr(err)
			}
			log.Println("=====================")
			spinner.Stop()

		},
	}
	cmd.PersistentFlags().StringVar(&action, "action", "", "Action for the platform that should be executed on current resource with params in a file")
	cmd.PersistentFlags().StringVar(&file, "file", "", "File containing the parameters for the action")
	cmd.PersistentFlags().BoolVar(&la, "la", false, "List down all the actions that can be performed on the current resource")

	return cmd
}

func getDataworkerAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.GetOps:
		return dpb.DataWorkerAgentActions_LIST_DATAWORKER.String()
	case pkg.DescribeOps:
		return dpb.DataWorkerAgentActions_DESCRIBE_DATAWORKER.String()
	case pkg.ActOps:
		return action
	default:
		return pkg.ErrInvalidAction.Error()
	}
}

// func dataWorkerActions(parentCmd string, args []string) {
// 	switch parentCmd {
// 	case pkg.GetOps:
// 		getDataWorker()
// 	case pkg.DescribeOps:
// 		describeDataWorker(args)
// 	default:
// 		cobra.CheckErr("Invalid action")
// 	}
// }

// func getDataWorker() {
// 	err := vapusGlobals.VapusPlatformClient.ListActions(dpb.DataWorkerAgentActions_LIST_DATAWORKER.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeDataWorker(args []string) {
// 	if len(args) < 1 {
// 		cobra.CheckErr("Invalid number of arguments, please provide the dataWorker ID")
// 	}
// 	err := vapusGlobals.VapusPlatformClient.DescribeActions(dpb.DataWorkerAgentActions_DESCRIBE_DATAWORKER.String(), viper.GetString(currentAccessToken), args[0])
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
