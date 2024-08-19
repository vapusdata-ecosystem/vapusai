package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
)

var (
	loginDataProduct string
)

// authCmd represents the auth command
func NewDataProductAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   dataProductAuthCmd,
		Short: "Login to the VapusData platform instance using Authenticator",
		Long:  `This command is used to login to the VapusData platform`,
		Run: func(cmd *cobra.Command, args []string) {
			generateDataProductAccessToken(args)
		},
	}
	cmd.Flags().StringVar(&loginDataProduct, "dataproduct", "", "uses provided data product context for logging in")
	return cmd
}

func generateDataProductAccessToken(args []string) {
	var err error
	if loginDataProduct == "" {
		cobra.CheckErr(pkg.ErrMissingDataProductLogin)
	}
	accessToken := viper.GetString(currentAccessToken)
	newDPAccessToken, err := vapusGlobals.VapusPlatformClient.RetrieveDataProductAccessToken(context.Background(), accessToken, loginDataProduct)
	if err != nil {
		vapusGlobals.logger.Error().Err(err).Msg("failed to retrieve platform access token")
		cobra.CheckErr(err)
	}

	viper.Set(currentProductAccessToken, newDPAccessToken)
	err = viper.WriteConfig()
	if err != nil {
		cobra.CheckErr(err)
	}
	vapusGlobals.logger.Info().Msgf("successfully logged in to data product - %v", loginDataProduct)
}
