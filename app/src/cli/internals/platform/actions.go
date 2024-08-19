package plclient

import (
	"context"
	"encoding/json"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
	dpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/domain/v1alpha1"
	mpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/platform/v1alpha1"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
	dmodels "github.com/vapusdata-ecosystem/vapusai-studio/core/models"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

func (x *VapusPlatformClient) ConfigureDomain(ctx context.Context, fileBytes []byte) error {
	spec := &mpb.Domain{}
	err := yaml.Unmarshal(fileBytes, spec)
	if err != nil {
		return err
	}
	response, err := x.DomainConn.DomainInterface(ctx, &dpb.DomainInterfaceRequest{
		Actions: []dpb.DomainAgentActions{dpb.DomainAgentActions_CONFIGURE_DOMAIN},
		Spec:    spec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Domain Configured Successfully with ID:", x.logger, response.Output.Domains[0].DomainId)
	return nil
}

func (x *VapusPlatformClient) ConfigureDataSource(ctx context.Context, fileBytes []byte) error {
	spec := &mpb.DataSource{}
	err := jsonpb.Unmarshal(fileBytes, spec)
	if err != nil {
		x.logger.Error().Msgf("Error in unmarshalling the file using JSONPB: %v", err)
		return err
	}
	response, err := x.DomainConn.DataSourceInterface(ctx, &dpb.DataSourceInterfaceRequest{
		Actions: []dpb.DataSourceAgentActions{dpb.DataSourceAgentActions_CONFIGURE_DATASOURCE},
		Spec:    spec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Data Source Configured Successfully with ID:", x.logger, response.Output.DataSources[0].DataSourceId)
	return nil
}

func (x *VapusPlatformClient) ConfigureDataWorker(ctx context.Context, fileBytes []byte) error {
	spec := &dmodels.DataWorker{}
	err := json.Unmarshal(fileBytes, spec)
	if err != nil {
		x.logger.Error().Msgf("Error in unmarshalling the file using JSONPB: %v", err)
		return err
	}
	val, ok := mpb.DataWorkerType_value[spec.DataWorkerType]
	if !ok {
		return pkg.ErrInvalidDataWorkerType
	}
	response, err := x.DomainConn.DataWorkerInterface(ctx, &dpb.DataWorkerInterfaceRequest{
		Actions:        []dpb.DataWorkerAgentActions{dpb.DataWorkerAgentActions_CONFIGURE_DATAWORKER},
		DataWorkerType: mpb.DataWorkerType(val),
		Spec:           string(fileBytes),
		Format:         mpb.FileFormats(mpb.FileFormats_value[x.inputFormat]),
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Data Worker Configured Successfully with ID:", x.logger, response.Output.DataWorkers[0].DataWorkerId)
	return nil
}

func (x *VapusPlatformClient) ConfigureDataProduct(ctx context.Context, fileBytes []byte) error {
	spec := &mpb.DataProduct{}
	err := jsonpb.Unmarshal(fileBytes, spec)
	if err != nil {
		x.logger.Error().Msgf("Error in unmarshalling the file using JSONPB: %v", err)
		return err
	}
	response, err := x.DomainConn.DataProductInterface(ctx, &dpb.DataProductInterfaceRequest{
		Actions: []dpb.DataProductAgentActions{dpb.DataProductAgentActions_CONFIGURE_DATAPRODUCT},
		Spec:    spec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Data Product Configured Successfully with ID:", x.logger, response.Output.DataProducts[0].DataProductId)
	return nil
}

func (x *VapusPlatformClient) BuildDataProduct(ctx context.Context, fileBytes []byte) error {
	spec := &mpb.DataProduct{}
	if len(x.ActionHandler.Identifier) < 1 {
		err := jsonpb.Unmarshal(fileBytes, spec)
		if err != nil {
			x.logger.Error().Msgf("Error in unmarshalling the file using JSONPB:", err)
			return err
		}
	} else {
		spec.DataProductId = x.ActionHandler.Identifier[0]
	}
	response, err := x.DomainConn.DataProductInterface(ctx, &dpb.DataProductInterfaceRequest{
		Actions: []dpb.DataProductAgentActions{dpb.DataProductAgentActions_BUILD_DATAPRODUCT},
		Spec:    spec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Data Product Built Successfully with ID:", x.logger, response.Output.DataProducts[0].DataProductId)
	return nil
}

func (x *VapusPlatformClient) PublishDataProduct(ctx context.Context, fileBytes []byte) error {
	spec := &mpb.DataProduct{}
	if len(x.ActionHandler.Identifier) < 1 {
		err := jsonpb.Unmarshal(fileBytes, spec)
		if err != nil {
			x.logger.Error().Msgf("Error in unmarshalling the file using JSONPB:", err)
			return err
		}
	} else {
		spec.DataProductId = x.ActionHandler.Identifier[0]
	}
	response, err := x.DomainConn.DataProductInterface(ctx, &dpb.DataProductInterfaceRequest{
		Actions: []dpb.DataProductAgentActions{dpb.DataProductAgentActions_PUBLISH_DATAPRODUCT},
		Spec:    spec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Data Product Published Successfully with ID: %v", x.logger, response.Output.DataProducts[0].DataProductId)
	return nil
}

func (x *VapusPlatformClient) DeployDataWorkers(ctx context.Context, fileBytes []byte) error {
	spec := &mpb.DataProduct{}
	if len(x.ActionHandler.Identifier) < 1 {
		err := jsonpb.Unmarshal(fileBytes, spec)
		if err != nil {
			x.logger.Error().Msgf("Error in unmarshalling the file using JSONPB:", err)
			return err
		}
	} else {
		spec.DataProductId = x.ActionHandler.Identifier[0]
	}
	response, err := x.DomainConn.DataProductInterface(ctx, &dpb.DataProductInterfaceRequest{
		Actions: []dpb.DataProductAgentActions{dpb.DataProductAgentActions_DEPLOY_WORKERS},
		Spec:    spec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Data Product workers are deployed successfully with ID:", x.logger, response.Output.DataProducts[0].DataProductId)
	return nil
}

func (x *VapusPlatformClient) DeployDataProduct(ctx context.Context, fileBytes []byte) error {
	spec := &mpb.DataProduct{}
	if len(x.ActionHandler.Identifier) < 1 {
		err := jsonpb.Unmarshal(fileBytes, spec)
		if err != nil {
			x.logger.Error().Msgf("Error in unmarshalling the file using JSONPB:", err)
			return err
		}
	} else {
		spec.DataProductId = x.ActionHandler.Identifier[0]
	}
	response, err := x.DomainConn.DataProductInterface(ctx, &dpb.DataProductInterfaceRequest{
		Actions: []dpb.DataProductAgentActions{dpb.DataProductAgentActions_DEPLOY_DATAPRODUCT},
		Spec:    spec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Data Product Deployed Successfully with ID:", x.logger, response.Output.DataProducts[0].DataProductId)
	return nil
}

func (x *VapusPlatformClient) SearchDataProducts(ctx context.Context) error {
	response, err := x.PlConn.VapusSearch(ctx, &pb.VapusSearchRequest{
		Q: x.ActionHandler.SearchQ,
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Data Product Search Result:", x.logger)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Name", "Data Catalog", "Data ProductId", "Owner Domain"})
	for _, item := range response.DataProductResults {
		tw.AppendRow(table.Row{item.DataProductName, item.Catalog, item.OwnerDomain, item.DataProductId})
		tw.AppendSeparator()
	}
	tw.Render()
	return nil
}

func (x *VapusPlatformClient) DescribeDataSourceMetaData(ctx context.Context) error {
	response, err := x.DomainConn.DataSourceInterface(ctx, &dpb.DataSourceInterfaceRequest{
		Actions: []dpb.DataSourceAgentActions{dpb.DataSourceAgentActions_LIST_DATASOURCE_METADATA},
		Query: &mpb.SearchParam{
			Search: []*mpb.MapList{
				{
					Key:    dmodels.DataSourceSK.String(),
					Values: x.ActionHandler.Identifier,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	if len(response.Output.GetMetadata()) < 1 || len(response.Output.GetDataSources()) < 1 {
		return pkg.ErrMetaData404
	}
	metaData := response.Output.GetMetadata()[0]
	dataSource := response.Output.GetDataSources()[0]
	pkg.LogTitlesf("MetaData for Data Source:", x.logger, response.Output.DataSources[0].DataSourceId)

	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Owner Domain", "service Provider", "Storage Engine", "Service Name", "status"})
	tw.AppendRow(table.Row{metaData.SourceMetaData.DataSource.Domain.Name,
		dataSource.Attributes.ServiceProvider,
		dataSource.Attributes.StorageEngine,
		dataSource.Attributes.ServiceName,
		dataSource.Status})
	tw.AppendSeparator()
	tw.Render()

	schemes := metaData.SourceMetaData.DataSource.DataSchema
	pkg.LogTitles("Description", x.logger)
	x.logger.Info().Msg(metaData.SourceMetaData.DataSource.Description)

	pkg.LogTitles("DataTables in Database", x.logger)
	tw1 := pkg.NewTableWritter()
	tw1.AppendHeader(table.Row{"Database", "Data Tables", "Data Fields", "Total Rows", "Created At", "Last Updated At", "Table Type", "Nature"})
	for _, item := range schemes.DataTables {
		fields := []interface{}{}
		for _, field := range item.Fields {
			// val := map[string]interface{}{
			// 	"Field": field.Field,
			// 	"Type":  field.Type,
			// 	"Null":  field.Null,
			// }
			fields = append(fields, field.Field)
		}
		tw1.AppendRow(table.Row{metaData.SourceMetaData.DataSource.DataSchema.DataSourceIdentifier,
			item.Name, pkg.NewListWritter(fields, list.StyleBulletSquare).Render(), item.TotalRows, pkg.UnixTransformer(int64(item.CreatedAt)), pkg.UnixTransformer(int64(item.LastUpdatedAt)), item.TableType, item.Nature.String()})
		tw1.AppendSeparator()
	}
	tw1.Render()
	pkg.LogTitles("Constraints in Database", x.logger)
	tw2 := pkg.NewTableWritter()
	tw2.AppendHeader(table.Row{"Contraint Name", "Contraint Type", "Data Table", "Enforced"})
	for _, item := range schemes.Constraints {
		tw2.AppendRow(table.Row{item.ConstraintName, item.ConstraintType, item.TableName, item.Enforced})
		tw2.AppendSeparator()
	}
	tw2.Render()

	pkg.LogTitles("Lineage", x.logger)
	tw3 := pkg.NewTableWritter()
	tw3.AppendHeader(table.Row{"Sync At", "Digest", "Is Latest", "Diff from last sync"})
	for _, item := range metaData.SourceMetaData.DataSource.Lineage {
		tw3.AppendRow(table.Row{pkg.UnixTransformer(item.SyncAt), item.Digest.Digest, item.IsLatest, item.MetadataDiff})
		tw3.AppendSeparator()
	}
	tw3.Render()
	pkg.LogTitles("Data Compliance", x.logger)
	tw4 := pkg.NewTableWritter()
	tw4.AppendHeader(table.Row{"Field Name", "Data Table", "Type", "Total Entries", "Null Entries", "Compliances", "Current First Record At", "Latest Record At"})
	for _, item := range dataSource.GetComplianceFields() {
		tw4.AppendRow(table.Row{item.Name, item.DataTable, item.FieldType, item.TotalEntries, item.TotalNullEntries, item.ComplianceTypes, item.CurrentFirstRecordAt, item.LatestRecordAt})
		tw4.AppendSeparator()
	}
	tw4.Render()
	return nil
}

func (x *VapusPlatformClient) ListDataSourceMetaData(ctx context.Context) error {
	response, err := x.DomainConn.DataSourceInterface(ctx, &dpb.DataSourceInterfaceRequest{
		Actions: []dpb.DataSourceAgentActions{dpb.DataSourceAgentActions_LIST_DATASOURCE_METADATA},
		Query: &mpb.SearchParam{
			Search: []*mpb.MapList{
				{
					Key:    dmodels.DataSourceSK.String(),
					Values: x.ActionHandler.Identifier,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	if len(response.Output.GetMetadata()) < 1 || len(response.Output.GetDataSources()) < 1 {
		return pkg.ErrMetaData404
	}
	metaData := response.Output.GetMetadata()[0]
	dataSource := response.Output.GetDataSources()[0]
	pkg.LogTitlesf("MetaData for Data Source:", x.logger, response.Output.DataSources[0].DataSourceId)

	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Owner Domain", "service Provider", "Storage Engine", "Service Name", "status"})
	tw.AppendRow(table.Row{metaData.SourceMetaData.DataSource.Domain.Name,
		dataSource.Attributes.ServiceProvider,
		dataSource.Attributes.StorageEngine,
		dataSource.Attributes.ServiceName,
		dataSource.Status})
	tw.AppendSeparator()
	tw.Render()

	schemes := metaData.SourceMetaData.DataSource.DataSchema
	pkg.LogTitles("Description", x.logger)
	x.logger.Info().Msg(metaData.SourceMetaData.DataSource.Description)

	pkg.LogTitles("DataTables in Database", x.logger)
	tw1 := pkg.NewTableWritter()
	tw1.AppendHeader(table.Row{"Database", "Data Tables", "Data Fields", "Total Rows", "Created At", "Last Updated At", "Table Type", "Nature"})
	for _, item := range schemes.DataTables {
		fields := []interface{}{}
		for _, field := range item.Fields {
			// val := map[string]interface{}{
			// 	"Field": field.Field,
			// 	"Type":  field.Type,
			// 	"Null":  field.Null,
			// }
			fields = append(fields, field.Field)
		}
		tw1.AppendRow(table.Row{metaData.SourceMetaData.DataSource.DataSchema.DataSourceIdentifier,
			item.Name, pkg.NewListWritter(fields, list.StyleBulletSquare).Render(), item.TotalRows, pkg.UnixTransformer(int64(item.CreatedAt)), pkg.UnixTransformer(int64(item.LastUpdatedAt)), item.TableType, item.Nature.String()})
		tw1.AppendSeparator()
	}
	tw1.Render()
	pkg.LogTitles("Constraints in Database", x.logger)
	tw2 := pkg.NewTableWritter()
	tw2.AppendHeader(table.Row{"Contraint Name", "Contraint Type", "Data Table", "Enforced"})
	for _, item := range schemes.Constraints {
		tw2.AppendRow(table.Row{item.ConstraintName, item.ConstraintType, item.TableName, item.Enforced})
		tw2.AppendSeparator()
	}
	tw2.Render()

	pkg.LogTitles("Lineage", x.logger)
	tw3 := pkg.NewTableWritter()
	tw3.AppendHeader(table.Row{"Sync At", "Digest", "Is Latest", "Diff from last sync"})
	for _, item := range metaData.SourceMetaData.DataSource.Lineage {
		tw3.AppendRow(table.Row{pkg.UnixTransformer(item.SyncAt), item.Digest.Digest, item.IsLatest, item.MetadataDiff})
		tw3.AppendSeparator()
	}
	tw3.Render()
	pkg.LogTitles("Data Compliance", x.logger)
	tw4 := pkg.NewTableWritter()
	tw4.AppendHeader(table.Row{"Field Name", "Data Table", "Type", "Total Entries", "Null Entries", "Compliances", "Current First Record At", "Latest Record At"})
	for _, item := range dataSource.GetComplianceFields() {
		tw4.AppendRow(table.Row{item.Name, item.DataTable, item.FieldType, item.TotalEntries, item.TotalNullEntries, item.ComplianceTypes, item.CurrentFirstRecordAt, item.LatestRecordAt})
		tw4.AppendSeparator()
	}
	tw4.Render()
	return nil
}
