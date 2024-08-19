package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	loginDomain string
)

// authCmd represents the auth command
func NewDomainAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   domainAuthCmd,
		Short: "Login to the VapusData platform instance using Authenticator",
		Long:  `This command is used to login to the VapusData platform`,
		Run: func(cmd *cobra.Command, args []string) {
			generateDomainAccessToken(args)
		},
	}

	cmd.Flags().StringVar(&loginDomain, "domain", "", "uses provided domain context for logging in")
	return cmd
}

func generateDomainAccessToken(args []string) {
	var err error
	if loginDomain == "" {
		vapusGlobals.logger.Info().Msg("no domain provided for login, system will login to default domain")
	}
	accessToken := viper.GetString(currentAccessToken)
	newAccessToken, err := vapusGlobals.VapusPlatformClient.RetrievePlatformAccessToken(context.Background(), accessToken, loginDomain)
	if err != nil {
		vapusGlobals.logger.Error().Err(err).Msg("failed to retrieve platform access token")
		cobra.CheckErr(err)
	}

	viper.Set(currentAccessToken, newAccessToken)
	err = viper.WriteConfig()
	if err != nil {
		vapusGlobals.logger.Error().Err(err).Msg("failed to write new access token to config")
		cobra.CheckErr(err)
	}
	vapusGlobals.logger.Info().Msgf("successfully logged in to domain - %v", loginDomain)
}
