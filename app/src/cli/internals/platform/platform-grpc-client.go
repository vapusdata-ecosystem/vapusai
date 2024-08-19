package plclient

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rs/zerolog"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"sigs.k8s.io/yaml"

	dpb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/domain/v1alpha1"
	pb "github.com/vapusdata-ecosystem/vapusai-studio/apis/protos/platform/v1alpha1"
	pkg "github.com/vapusdata-ecosystem/vapusai-studio/cli/pkgs"
	"github.com/vapusdata-ecosystem/vapusai-studio/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusai-studio/core/pkgs/utils"
)

var AgentGoals = map[string][]interface{}{
	"account":        getAgentOps(pb.AccountAgentActions_name),
	"datamesh":       getAgentOps(pb.DataMeshAgentActions_name),
	"vapussearch":    getAgentOps(pb.VapusSearchType_name),
	"datacatalog":    getAgentOps(pb.DataCatalogAgentActions_name),
	"datacompliance": getAgentOps(nil),
	"domain":         getAgentOps(dpb.DomainAgentActions_name),
	"datasource":     getAgentOps(dpb.DataSourceAgentActions_name),
	"dataworker":     getAgentOps(dpb.DataWorkerAgentActions_name),
	"dataproduct":    getAgentOps(dpb.DataProductAgentActions_name),
	"user":           getAgentOps(pb.UserAgentOperations_name),
	"authorization":  getAgentOps(pb.AccessTokenAgentUtility_name),
	"utility":        getAgentOps(nil),
	"observability":  getAgentOps(dpb.ObservabilityAgentReflexes_name),
}

type ActionHandlerOpts struct {
	ParentCmd   string
	Args        []string
	Identifier  []string
	Action      string
	File        string
	La          bool
	AccessToken string
	SearchQ     string
}

type VapusPlatformClient struct {
	Host              string
	PlConn            pb.VapusDataPlatformServiceClient
	UserConn          pb.VapusDataPlatformUserServiceClient
	DomainConn        dpb.DomainServiceClient
	grpcClient        *pbtools.GrpcClient
	CaCertFile        string
	ClientCertFile    string
	ClientKeyFile     string
	ValidTill         time.Time
	Error             error
	logger            zerolog.Logger
	ResourceActionMap map[string][]interface{}
	AccessToken       string
	marshaller        func(m protoreflect.ProtoMessage) ([]byte, error)
	inputFormat       string
	ActionHandler     ActionHandlerOpts
}

func getAgentOps(enum_map map[int32]string) []interface{} {
	var ops []interface{}
	if enum_map == nil {
		return ops
	}
	for _, v := range enum_map {
		ops = append(ops, v)
	}
	return ops
}

func NewPlatFormClient(params map[string]string, logger zerolog.Logger) (*VapusPlatformClient, error) {
	url, ok := params["url"]
	if !ok {
		return nil, errors.New("url is required, missing from the context")
	}
	// namespace, ok := params["namespace"]
	// if !ok {
	// 	return nil, errors.New("namespace is required, missing from the context")
	// }
	port, ok := params["port"]
	if !ok {
		return nil, errors.New("port is required, missing from the context")
	}
	portI, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.Join(err, errors.New("port is not a valid integer"))
	}
	dns := fmt.Sprintf("%s:%d", url, portI)
	// dns = "localhost:9013"

	telnet, err := net.DialTimeout("tcp", dns, 1*time.Second)
	if err != nil {
		return nil, err
	}
	defer telnet.Close()

	grpcClient := pbtools.NewGrpcClient(logger,
		pbtools.ClientWithInsecure(true),
		pbtools.ClientWithServiceAddress(dns))
	cl := &VapusPlatformClient{
		Host:       dns,
		PlConn:     pb.NewVapusDataPlatformServiceClient(grpcClient.Connection),
		grpcClient: grpcClient,
		UserConn:   pb.NewVapusDataPlatformUserServiceClient(grpcClient.Connection),
		DomainConn: dpb.NewDomainServiceClient(grpcClient.Connection),
		logger:     logger,
		marshaller: jsonpb.MarshalOptions{
			UseProtoNames:   true,
			UseEnumNumbers:  false,
			Indent:          "  ",
			EmitUnpopulated: true,
		}.Marshal,
		ActionHandler: ActionHandlerOpts{},
	}
	return cl, nil
}

func (x *VapusPlatformClient) ListResourceActions(resource string) {
	xx := text.FormatUpper.Apply("Actions for Domain Resource: ")
	xx = text.Underline.Sprintf(xx)
	x.logger.Info().Msgf("\n%v", xx)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Actions", "Version", "Commands"})
	for _, action := range x.ResourceActionMap[resource] {
		if strings.Contains(action.(string), "INVALID") {
			continue
		}
		tw.AppendRow(table.Row{action, "v1alpha1", pkg.APPNAME + " act " + resource + " --file <Input File> --action " + action.(string)})
		tw.AppendSeparator()
	}
	tw.Render()
}

func (x *VapusPlatformClient) HandleAction() error {
	if x.ActionHandler.ParentCmd == "" {
		return errors.New("invalid operations")
	}
	x.logger.Info().Msgf("Handling Action: %s", x.ActionHandler.Action)
	switch x.ActionHandler.ParentCmd {
	case pkg.GetOps:
		return x.GetActions()
	case pkg.DescribeOps:
		return x.DescribeActions()
	case pkg.ActOps:
		return x.PerformAct()
	case pkg.SearchOpts:
		return x.Search()
	default:
		return pkg.ErrInvalidAction
	}
}

func (x *VapusPlatformClient) GetActions() error {
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	switch x.ActionHandler.Action {
	case pb.AccountAgentActions_LIST_ACCOUNT.String():
		return x.ListAccount(newCtx)
	case pb.UserAgentOperations_GET_USER.String():
		return x.ListUser(newCtx)
	case dpb.DomainAgentActions_LIST_DOMAINS.String():
		return x.ListDomains(newCtx)
	case dpb.DataSourceAgentActions_LIST_DATASOURCE.String():
		return x.ListDataSources(newCtx)
	case pb.DataMeshAgentActions_LIST_DATAMESH.String():
		return x.ListDataMesh(newCtx)
	case dpb.DataWorkerAgentActions_LIST_DATAWORKER.String():
		return x.ListDataWorkers(newCtx)
	case dpb.DataProductAgentActions_LIST_DATAPRODUCT.String():
		return x.ListDataProducts(newCtx)
	case pb.DataCatalogAgentActions_LIST_CATALOG.String():
		return x.ListDataCatalogs(newCtx)
	default:
		return pkg.ErrInvalidAction
	}
}

func (x *VapusPlatformClient) DescribeActions() error {
	if len(x.ActionHandler.Args) < 1 {
		return pkg.ErrNoArgs
	}
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	switch x.ActionHandler.Action {
	case pb.AccountAgentActions_LIST_ACCOUNT.String():
		return x.DescribeAccount(newCtx)
	case pb.UserAgentOperations_GET_USER.String():
		return x.DescribeUser(newCtx)
	case dpb.DomainAgentActions_LIST_DOMAINS.String():
		return x.DescribeDomains(newCtx)
	case dpb.DataSourceAgentActions_DESCRIBE_DATASOURCE.String():
		return x.DescribeDataSources(newCtx, x.ActionHandler.Args[0])
	case pb.DataMeshAgentActions_LIST_DATAMESH.String():
		return x.DescribeDataMesh(newCtx)
	case dpb.DataWorkerAgentActions_LIST_DATAWORKER.String():
		return x.DescribeDataWorkers(newCtx, x.ActionHandler.Args[0])
	case dpb.DataProductAgentActions_LIST_DATAPRODUCT.String():
		return x.DescribeDataProducts(newCtx, x.ActionHandler.Args[0])
	case pb.DataCatalogAgentActions_LIST_CATALOG.String():
		return x.DescribeDataCatalogs(newCtx)
	default:
		return pkg.ErrInvalidAction
	}
}

func (x *VapusPlatformClient) PerformAct() error {
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	var fileBytes []byte
	if x.ActionHandler.File != "" {
		bbytes, err := dmutils.ReadFile(x.ActionHandler.File)
		if err != nil {
			return err
		}
		fileBytes, err = yaml.YAMLToJSON(bbytes)
		if err != nil {
			return err
		}
		x.inputFormat = strings.ToUpper(dmutils.GetConfFileType(x.ActionHandler.File))
	}
	if len(x.ActionHandler.Identifier) < 1 && x.ActionHandler.File == "" {
		return pkg.ErrNoArgs
	}
	switch x.ActionHandler.Action {
	case dpb.DomainAgentActions_CONFIGURE_DOMAIN.String():
		return x.ConfigureDomain(newCtx, fileBytes)
	case dpb.DataSourceAgentActions_CONFIGURE_DATASOURCE.String():
		return x.ConfigureDataSource(newCtx, fileBytes)
	case dpb.DataWorkerAgentActions_CONFIGURE_DATAWORKER.String():
		return x.ConfigureDataWorker(newCtx, fileBytes)
	case dpb.DataProductAgentActions_CONFIGURE_DATAPRODUCT.String():
		return x.ConfigureDataProduct(newCtx, fileBytes)
	case dpb.DataProductAgentActions_BUILD_DATAPRODUCT.String():
		return x.BuildDataProduct(newCtx, fileBytes)
	case dpb.DataProductAgentActions_PUBLISH_DATAPRODUCT.String():
		return x.PublishDataProduct(newCtx, fileBytes)
	case dpb.DataProductAgentActions_DEPLOY_WORKERS.String():
		return x.DeployDataWorkers(newCtx, fileBytes)
	case dpb.DataProductAgentActions_DEPLOY_DATAPRODUCT.String():
		return x.DeployDataProduct(newCtx, fileBytes)
	case dpb.DataSourceAgentActions_LIST_DATASOURCE_METADATA.String():
		return x.ListDataSourceMetaData(newCtx)
	case dpb.DataSourceAgentActions_DESCRIBE_DATASOURCE_METADATA.String():
		return x.DescribeDataSourceMetaData(newCtx)
	default:
		return pkg.ErrInvalidAction
	}
}

func (x *VapusPlatformClient) Search() error {
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	switch x.ActionHandler.Action {
	case pb.VapusSearchType_DATAPRODUCTS.String():
		return x.SearchDataProducts(newCtx)
	default:
		return pkg.ErrInvalidAction
	}
}

func (x *VapusPlatformClient) Close() {
	x.PlConn = nil
	x.UserConn = nil
	x.DomainConn = nil
	x.grpcClient.Close()
}

func (x *VapusPlatformClient) PrintDescribe(data protoreflect.ProtoMessage, resource string) {
	jsonData, err := x.marshaller(data)
	if err != nil {
		x.logger.Error().Msgf("Error in marshaling %v details", resource)
	}
	bytes, err := yaml.JSONToYAML([]byte(jsonData))
	if err != nil {
		x.logger.Error().Msgf("Error in formatting %v details", resource)
	}
	x.logger.Info().Msgf("\n%s", string(bytes))
}
