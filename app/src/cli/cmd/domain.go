package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	dpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/domain/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai-studio/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
)

func NewDomainCmd() *cobra.Command {
	var action, file string
	var la bool
	cmd := &cobra.Command{
		Use:   domainResource,
		Short: "This command is interface to interact with the platform for domain resources.",
		Long:  `This command is interface to interact with the platform for domain resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			defer vapusGlobals.VapusPlatformClient.Close()
			if la {
				vapusGlobals.VapusPlatformClient.ListResourceActions("domain")
				return
			}
			spinner := pkg.GetSpinner(36)
			spinner.Prefix = " Creating domain in the platform"
			spinner.Start()
			vapusGlobals.VapusPlatformClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(currentAccessToken),
				Action:      getDomainAction(cmd.Parent().Use, action),
				File:        file,
			}
			err := vapusGlobals.VapusPlatformClient.HandleAction()
			if err != nil {
				spinner.Stop()
				cobra.CheckErr(err)
			}
			spinner.Stop()
		},
	}
	cmd.PersistentFlags().StringVar(&action, "action", "", "Action for the platform that should be executed on current resource with params in a file")
	cmd.PersistentFlags().StringVar(&file, "file", "", "File containing the parameters for the action")
	cmd.PersistentFlags().BoolVar(&la, "la", false, "List down all the actions that can be performed on the domain resource")
	return cmd
}

func getDomainAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.GetOps:
		return dpb.DomainAgentActions_LIST_DOMAINS.String()
	case pkg.DescribeOps:
		return dpb.DomainAgentActions_LIST_DOMAINS.String()
	case pkg.ActOps:
		return action
	default:
		return pkg.ErrInvalidAction.Error()
	}
}

// func (x DomainHandler) getDomain() {
// 	err := vapusGlobals.VapusPlatformClient.ListActions(dpb.DomainAgentActions_LIST_DOMAINS.String(), x.accessToken)
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func (x DomainHandler) describeDomain() {
// 	if len(x.args) < 1 {
// 		cobra.CheckErr("Invalid number of arguments, please provide the domain ID")
// 	}
// 	err := vapusGlobals.VapusPlatformClient.DescribeActions(dpb.DomainAgentActions_LIST_DOMAINS.String(), x.accessToken, x.args[0])
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func (x DomainHandler) act() {
// 	if x.action == "" {
// 		cobra.CheckErr("No action provided")
// 	}
// 	if x.file == "" {
// 		cobra.CheckErr("No input provided")
// 	}
// 	err := vapusGlobals.VapusPlatformClient.PerformAct(x.action, x.accessToken, x.file)
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
