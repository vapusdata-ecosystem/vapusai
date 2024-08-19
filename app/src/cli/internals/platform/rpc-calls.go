package plclient

import (
	"context"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
	dpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/domain/v1alpha1"
	mpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/platform/v1alpha1"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
	dmodels "github.com/vapusdata-ecosystem/vapusai-studio/core/models"
	"google.golang.org/grpc/metadata"
)

func (x *VapusPlatformClient) RetrieveLoginURL() (*pb.LoginHandlerResponse, error) {
	return x.UserConn.LoginHandler(context.Background(), &mpb.EmptyRequest{})
}

func (x *VapusPlatformClient) RetrieveAccessToken(code, callbackURL string) (string, string, error) {
	result, err := x.UserConn.LoginCallback(context.Background(), &pb.LoginCallBackRequest{Code: code, CallbackURL: callbackURL})
	if err != nil {
		return "", "", err
	}
	return result.Token.GetAccessToken(), result.Token.GetIdToken(), nil
}

func (x *VapusPlatformClient) RetrievePlatformAccessToken(ctx context.Context, token, domain string) (string, error) {
	result, err := x.UserConn.AccessTokenInterface(pkg.GetBearerCtx(ctx, token), &pb.AccessTokenInterfaceRequest{Domain: domain, Utility: pb.AccessTokenAgentUtility_DOMAIN_LOGIN})
	if err != nil {
		return "", err
	}
	return result.Token.AccessToken, nil
}

func (x *VapusPlatformClient) RetrieveDataProductAccessToken(ctx context.Context, token, dataproductId string) (string, error) {
	result, err := x.UserConn.AccessTokenInterface(pkg.GetBearerCtx(ctx, token), &pb.AccessTokenInterfaceRequest{DataProduct: dataproductId, Utility: pb.AccessTokenAgentUtility_DATAPRODUCT_LOGIN})
	if err != nil {
		return "", err
	}
	return result.Token.AccessToken, nil
}

func (x *VapusPlatformClient) ListAccount(ctx context.Context) error {
	result, err := x.PlConn.AccountInterface(ctx, &pb.AccountInterfaceRequest{Actions: pb.AccountAgentActions_LIST_ACCOUNT})
	if err != nil {
		return err
	}
	pkg.LogTitles("Account Info of current loggedin instance: ", x.logger)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Account name", "Account Id", "Data Store", "Secret Store", "Artifact Store", "Authn Method", "Authz"})
	tw.AppendRow(table.Row{result.Output.Name, result.Output.AccountId, result.Output.BackendDataStorage.BesService, result.Output.BackendSecretStorage.BesService,
		result.Output.ArtifactStorage.BesService, result.Output.AuthnMethod, result.Output.DmAccessJWTKeys.SigningAlgorithm})
	tw.AppendSeparator()
	tw.Render()
	return nil
}

func (x *VapusPlatformClient) DescribeAccount(ctx context.Context) error {
	result, err := x.PlConn.AccountInterface(ctx, &pb.AccountInterfaceRequest{Actions: pb.AccountAgentActions_LIST_ACCOUNT})
	if err != nil {
		return err
	}

	pkg.LogTitles("Account Info: ", x.logger)
	x.PrintDescribe(result.Output, "account")
	return nil
}

func (x *VapusPlatformClient) ListDataMesh(ctx context.Context) error {
	result, err := x.PlConn.DataMeshInterface(ctx, &pb.DataMeshInterfaceRequest{
		Actions: []pb.DataMeshAgentActions{pb.DataMeshAgentActions_LIST_DATAMESH},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("List of Datamesh: ", x.logger)
	if len(result.Output.GetDatamesh()) == 0 {
		pkg.LogTitles("\nNo DataMesh found", x.logger)
	} else {
		tw := pkg.NewTableWritter()
		tw.AppendHeader(table.Row{"Name", "MeshId", "Status"})
		for _, dm := range result.Output.GetDatamesh() {
			tw.AppendRow(table.Row{dm.Name, dm.MeshId, dm.Status})
			tw.AppendSeparator()
		}
		tw.Render()
	}
	return nil
}

func (x *VapusPlatformClient) DescribeDataMesh(ctx context.Context) error {
	result, err := x.PlConn.DataMeshInterface(ctx, &pb.DataMeshInterfaceRequest{
		Actions: []pb.DataMeshAgentActions{pb.DataMeshAgentActions_LIST_DATAMESH},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Datamesh Info: ", x.logger)
	x.PrintDescribe(result.Output, "datamesh")
	return nil
}

func (x *VapusPlatformClient) ListUser(ctx context.Context) error {
	result, err := x.UserConn.GetUser(ctx, &mpb.EmptyRequest{})
	if err != nil {
		return err
	}
	pkg.LogTitles("Logged-in User Info: ", x.logger)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"UserId", "Name", "Added On"})
	tw1 := pkg.NewTableWritter()
	tw1.AppendHeader(table.Row{"Domain", "Roles"})
	for _, domain := range result.Output.GetDomainRoles() {
		tw1.AppendRow(table.Row{domain.DomainId, strings.Join(domain.Role, ",")})
	}
	tw1.Render()
	tw.AppendRow(table.Row{result.Output.UserId, result.Output.DisplayName, result.Output.InvitedOn})
	tw.AppendSeparator()
	tw.Render()
	return nil
}

func (x *VapusPlatformClient) DescribeUser(ctx context.Context) error {
	result, err := x.UserConn.GetUser(ctx, &mpb.EmptyRequest{})
	if err != nil {
		return err
	}
	pkg.LogTitles("Logged-in User Info: ", x.logger)
	x.PrintDescribe(result.Output, "user")
	return nil
}

func (x *VapusPlatformClient) ListDomains(ctx context.Context) error {
	result, err := x.DomainConn.DomainInterface(ctx, &dpb.DomainInterfaceRequest{
		Actions: []dpb.DomainAgentActions{dpb.DomainAgentActions_LIST_DOMAINS},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Logged-In user's Domain Info: ", x.logger)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Name", "Id", "Datasource Count", "status"})
	for _, dm := range result.Output.Domains {
		tw.AppendRow(table.Row{dm.Name, dm.DomainId, len(dm.DataSources), dm.Status})
		tw.AppendSeparator()
	}
	tw.Render()
	return nil
}

func (x *VapusPlatformClient) DescribeDomains(ctx context.Context) error {
	result, err := x.DomainConn.DomainInterface(ctx, &dpb.DomainInterfaceRequest{
		Actions: []dpb.DomainAgentActions{dpb.DomainAgentActions_LIST_DOMAINS},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Logged-In user's Domain Description: ", x.logger)
	x.PrintDescribe(result.Output.Domains[0], "domain")
	return nil
}

func (x *VapusPlatformClient) ListDataSources(ctx context.Context) error {
	result, err := x.DomainConn.DataSourceInterface(ctx, &dpb.DataSourceInterfaceRequest{
		Actions: []dpb.DataSourceAgentActions{dpb.DataSourceAgentActions_LIST_DATASOURCE},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("List of dataSources registered in current domain: ", x.logger)
	if len(result.Output.GetDataSources()) == 0 {
		pkg.LogTitles("\nNo Data Sources found", x.logger)
	} else {
		tw := pkg.NewTableWritter()
		tw.AppendHeader(table.Row{"Name", "Id", "service Provider", "Storage Engine", "Service Name", "status"})
		for _, source := range result.Output.GetDataSources() {
			tw.AppendRow(table.Row{source.Name, source.DataSourceId, source.Attributes.ServiceProvider, source.Attributes.StorageEngine, source.Attributes.ServiceName, source.Status})
			tw.AppendSeparator()
		}
		tw.Render()
	}
	return nil
}

func (x *VapusPlatformClient) DescribeDataSources(ctx context.Context, iden string) error {
	result, err := x.DomainConn.DataSourceInterface(ctx, &dpb.DataSourceInterfaceRequest{
		Actions: []dpb.DataSourceAgentActions{dpb.DataSourceAgentActions_DESCRIBE_DATASOURCE},
		Spec: &mpb.DataSource{
			DataSourceId: iden,
		},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Data Source Info: ", x.logger)
	x.PrintDescribe(result.Output.DataSources[0], "datasource")
	return nil
}

func (x *VapusPlatformClient) ListDataCatalogs(ctx context.Context) error {
	result, err := x.PlConn.DataCatalogInterface(ctx, &pb.DataCatalogInterfaceRequest{
		Actions: []pb.DataCatalogAgentActions{pb.DataCatalogAgentActions_GET_SELF_DATA_CATALOG},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("List of Data catalogs in the platforms: ", x.logger)
	if len(result.Output.GetDataCatalogs()) == 0 {
		pkg.LogTitles("\nNo Data catalogs found", x.logger)
	} else {
		tw := pkg.NewTableWritter()
		tw.AppendHeader(table.Row{"Name", "Catalog Id", "Data Mesh", "Data Product Count"})
		for _, dm := range result.Output.GetDataCatalogs() {
			tw.AppendRow(table.Row{dm.Name, dm.DataCatalogId, dm.MeshId, len(dm.DataProducts)})
			tw.AppendSeparator()
		}
		tw.Render()
	}
	return nil
}

func (x *VapusPlatformClient) DescribeDataCatalogs(ctx context.Context) error {
	result, err := x.PlConn.DataCatalogInterface(ctx, &pb.DataCatalogInterfaceRequest{
		Actions: []pb.DataCatalogAgentActions{pb.DataCatalogAgentActions_LIST_CATALOG},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Data Catalog Info: ", x.logger)
	x.PrintDescribe(result.Output.DataCatalogs[0], "datacatalog")
	return nil
}

func (x *VapusPlatformClient) ListDataProducts(ctx context.Context) error {
	result, err := x.DomainConn.DataProductInterface(ctx, &dpb.DataProductInterfaceRequest{
		Actions: []dpb.DataProductAgentActions{dpb.DataProductAgentActions_LIST_DATAPRODUCT},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("List of Data products in the current domain: ", x.logger)
	if len(result.Output.GetDataProducts()) == 0 {
		pkg.LogTitles("\nNo Data Products found", x.logger)
	} else {
		tw := pkg.NewTableWritter()
		tw.AppendHeader(table.Row{"Name", "Data Product id", "Data Catalog", "Status"})
		for _, dm := range result.Output.GetDataProducts() {
			tw.AppendRow(table.Row{dm.Name, dm.DataProductId, dm.DataCatalogId, dm.GetStatus()})
			tw.AppendSeparator()
		}
		tw.Render()
	}
	return nil
}

func (x *VapusPlatformClient) DescribeDataProducts(ctx context.Context, iden string) error {
	result, err := x.DomainConn.DataProductInterface(ctx, &dpb.DataProductInterfaceRequest{
		Actions: []dpb.DataProductAgentActions{dpb.DataProductAgentActions_DESCRIBE_DATAPRODUCT},
		Spec: &mpb.DataProduct{
			DataProductId: iden,
		},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("List of Data products in the current domain: ", x.logger)
	x.PrintDescribe(result.Output.GetDataProducts()[0], "dataproduct")
	return nil
}

func (x *VapusPlatformClient) ListDataWorkers(ctx context.Context) error {
	result, err := x.DomainConn.DataWorkerInterface(ctx, &dpb.DataWorkerInterfaceRequest{
		Actions: []dpb.DataWorkerAgentActions{dpb.DataWorkerAgentActions_LIST_DATAWORKER},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("List of Data workers in the current domain: ", x.logger)
	if len(result.Output.GetDataWorkers()) == 0 {
		pkg.LogTitles("\nNo Data Worker found", x.logger)
	} else {
		tw := pkg.NewTableWritter()
		tw.AppendHeader(table.Row{"Domain", "Data Worker Id", "Data Worker Type", "Status", "Created By"})
		for _, dm := range result.Output.GetDataWorkers() {
			tw.AppendRow(table.Row{dm.Domain, dm.GetDataWorkerId(), dm.GetDataWorkerType(), dm.Status, dm.GetAttributes().GetCreatedBy()})
			tw.AppendSeparator()
		}
		tw.Render()
	}
	return nil
}

func (x *VapusPlatformClient) DescribeDataWorkers(ctx context.Context, iden string) error {
	result, err := x.DomainConn.DataWorkerInterface(ctx, &dpb.DataWorkerInterfaceRequest{
		Actions: []dpb.DataWorkerAgentActions{dpb.DataWorkerAgentActions_DESCRIBE_DATAWORKER},
		Query: &mpb.SearchParam{
			Search: []*mpb.MapList{
				{
					Key:    dmodels.DataWorkerSK.String(),
					Values: []string{iden},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("List of Data workers in the current domain: ", x.logger)
	x.PrintDescribe(result.Output.GetDataWorkers()[0], "dataworker")
	return nil
}

func (x *VapusPlatformClient) ListPlatformSpec() {
	fm := []interface{}{}
	for _, f := range mpb.FileFormats_name {
		fm = append(fm, f)
	}
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Resource Name", "Spec Available", "Formats Available", "Generate command"})
	for _, spec := range mpb.RequestObjects_name {
		if spec == mpb.RequestObjects_INVALID_REQUEST_OBJECT.String() {
			continue
		}
		tw.AppendRow(table.Row{spec, true, pkg.NewListWritter(fm, list.StyleBulletSquare).Render(), "[PARENT_CMD] spec --name " + spec + " --generate-file=true --format yaml --with-fake=true"})
		tw.AppendSeparator()
	}
	tw.Render()
}

func (x *VapusPlatformClient) GeneratePlatformSpec(token, specName, format string, withFakeData bool) error {
	ctx := context.Background()
	md := metadata.Pairs("authorization", "Bearer "+token)
	newCtx := metadata.NewOutgoingContext(ctx, md)
	result, err := x.PlConn.GetSampleResourceConfiguration(newCtx, &pb.SampleResourceConfigurationOptions{
		Format:           mpb.FileFormats(mpb.FileFormats_value[format]),
		RequestObj:       mpb.RequestObjects(mpb.RequestObjects_value[specName]),
		PopulateFakeData: withFakeData,
	})
	if err != nil {
		return err
	}

	for _, f := range result.Output {
		fileName := strings.ToLower(f.GetRequestObj().String() + "." + strings.ToLower(format))
		x.logger.Info().Msgf("Sample %v spec generated with file name - %v \n", f.GetRequestObj(), fileName)
		err := os.WriteFile(fileName, []byte(f.GetFileContent()), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
