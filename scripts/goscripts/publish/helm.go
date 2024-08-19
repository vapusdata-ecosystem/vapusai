package vapuspublish

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/registry"
)

var (
	logger zerolog.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	helmRegistry, helmPassword,
	helmUsername, platformSvcDigest, platformSvcTag,
	vapusctlSvcDigest, vapusoperatorSvcDigest,
	vapusctlSvcTag, vapusoperatorSvcTag,
	etlworkerSvcDigest, etlworkerSvcTag,
	vapusDcSvcDigest, vapusDcSvcTag, appVersion string
	upload, bumpVersion bool
)

func init() {
	flag.StringVar(&helmRegistry, "helm-registry", "", "URL of registry")
	flag.StringVar(&helmUsername, "helm-registry-username", "", "URL of registry")
	flag.StringVar(&helmPassword, "helm-registry-password", "", "URL of registry")
	flag.StringVar(&platformSvcDigest, "platform-svc-digest", "", "Platform service digest")
	flag.StringVar(&platformSvcTag, "platform-svc-tag", "", "Platform service digest")
	flag.StringVar(&vapusctlSvcDigest, "vapusctl-svc-digest", "", "vapusctl service digest")
	flag.StringVar(&vapusoperatorSvcDigest, "vapusoperator-svc-digest", "", "vapusoperator service digest")
	flag.StringVar(&vapusctlSvcTag, "vapusctl-svc-tag", "", "vapusctl service digest")
	flag.StringVar(&vapusoperatorSvcTag, "vapusoperator-svc-tag", "", "vapusoperator service digest")
	flag.StringVar(&etlworkerSvcDigest, "etlworker-svc-digest", "", "etlworker service digest")
	flag.StringVar(&etlworkerSvcTag, "etlworker-svc-tag", "", "etlworker service digest")
	flag.StringVar(&vapusDcSvcDigest, "vapus-dc-svc-digest", "", "vapus-dc service digest")
	flag.StringVar(&vapusDcSvcTag, "vapus-dc-svc-tag", "", "vapus-dc service digest")
	flag.StringVar(&appVersion, "appVersion", "", "App version of the chart")
	flag.BoolVar(&upload, "upload", false, "Flag to control the upload of helm chart")
	flag.BoolVar(&bumpVersion, "bump-version", true, "Flag to bump the version of helm chart")
	flag.Parse()
}

type ArtifactValues struct {
	Tag    string `yaml:"tag"`
	Digest string `yaml:"digest"`
}

func PushOciImages() string {
	helmLocal := "../../deployments/helm-chart"
	helmOCI := "oci://" + helmRegistry
	chartName := "vapusdata-platform"
	tempDir, err := os.MkdirTemp("", "helm-chart-")
	logger.Info().Msgf("helmRegistry: %v", helmRegistry)
	logger.Info().Msgf("helmUsername: %v", helmUsername)
	logger.Info().Msgf("helmPassword: %v", helmPassword)
	if err != nil {
		logger.Info().Msgf("Error creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	chartFilePath := filepath.Join(helmLocal, chartName, "Chart.yaml")

	// Load the chart from the local directory
	chart, err := loader.LoadDir(filepath.Join(helmLocal, chartName))
	if err != nil {
		logger.Info().Msgf("Failed to load chart: %v", err)
	}

	// Package the chart
	chart.Metadata.Version = bumpChartVersion(chart.Metadata.Version)
	if bumpVersion {
		err = chartutil.SaveChartfile(chartFilePath, chart.Metadata)
		if err != nil {
			logger.Fatal().Msgf("Failed to save chart file: %v", err)
			return ""
		}
		logger.Info().Msgf("Bumped chart version to: %v", chart.Metadata.Version)
		return chart.Metadata.Version
	}

	if appVersion != "" {
		chart.Metadata.AppVersion = appVersion
	}
	// Save the chart file
	err = chartutil.SaveChartfile(chartFilePath, chart.Metadata)
	if err != nil {
		logger.Fatal().Msgf("Failed to save chart file: %v", err)
		return ""
	}

	// Update the values.yaml file
	err = updateVapusDataValues(chart, filepath.Join(helmLocal, chartName, "values.yaml"))
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to update values.yaml: %v", err)
		return ""
	}
	if upload {
		// Save the chart to a temporary directory
		chartPackagePath, err := chartutil.Save(chart, tempDir)
		if err != nil {
			logger.Fatal().Msgf("Failed to save chart: %v", err)
			return ""
		}

		// Initialize the registry client
		registryClient, err := registry.NewClient(
			registry.ClientOptDebug(true),
			registry.ClientOptWriter(os.Stdout),
		)
		if err != nil {
			logger.Fatal().Msgf("Failed to create registry client: %v", err)
			return ""
		}
		logger.Info().Msgf("Registry: %v", helmRegistry)
		// Login to the registry
		err = exec.Command("helm", "registry", "login", helmRegistry, "-u", helmUsername, "-p", helmPassword).Run()
		if err != nil {
			logger.Fatal().Msgf("Failed to login to registry using CLI: %v", err)
			return ""
		}

		defer func() {
			err = exec.Command("helm", "registry", "logout", helmRegistry).Run()
			if err != nil {
				logger.Info().Msgf("Failed to logout from registry using CLI: %v", err)
			}
		}()

		var output []byte
		output, err = exec.Command("helm", "push", chartPackagePath, helmOCI).CombinedOutput()
		if err != nil {
			logger.Fatal().Err(err).Msgf("Failed to push chart to OCI registry using CLI: %v", err)
			return ""

		}
		logger.Info().Msgf("Pushed chart to OCI registry using CLI: %v", string(output))
		// ociCredentials := registry.LoginOptBasicAuth(helmUsername, helmPassword)
		// registryClient.Login(helmRegistry, ociCredentials)
		defer registryClient.Logout(helmRegistry)
		// Push the chart package to the OCI registry with provenance data
		// result, err := registryClient.Push(chartPackagePath, helmOCI, registry.PushOptProvData(map[string]interface{}{}))
		// logger.Info().Msgf("Pushing chart to OCI registry: %v", helmOCI)
		// logger.Info().Msgf("Chart package path: %v", chartPackagePath)
		// result, err := registryClient.Push([]byte(chartPackagePath), helmOCI)
		// if err != nil {
		// 	logger.Info().Msgf("Failed to push chart to OCI registry: %v", err)
		// }
		digest := getDigestFromCmOp(string(output))
		logger.Info().Msgf("Pushed chart to OCI registry with digest: %v", digest)
		return chart.Metadata.Version
	} else {
		logger.Info().Msg("Skipping upload of helm chart")
		return chart.Metadata.Version
	}
}

func getDigestFromCmOp(text string) string {
	a := strings.Split(text, "Digest:")
	return strings.Trim(strings.Split(a[1], "\n")[0], " ")
}

func bumpChartVersion(current string) string {
	ver := strings.Split(current, ".")
	t, err := strconv.Atoi(ver[2])
	if err != nil {
		logger.Info().Msgf("Failed to convert string to int: %v", err)
		return current
	}
	if t < 200 {
		t += 1
		ver[2] = strconv.Itoa(t)
		return strings.Join(ver, ".")
	} else {
		t1, err := strconv.Atoi(ver[1])
		if err != nil {
			logger.Info().Msgf("Failed to convert string to int: %v", err)
			return current
		}
		t1 += 1
		ver[1] = strconv.Itoa(t1)
		ver[2] = "0"
		return strings.Join(ver, ".")
	}
}

func updateVapusDataValues(chart *chart.Chart, file string) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		logger.Info().Msgf("Failed to read file: %v", err)
		return err
	}
	result, err := chartutil.ReadValues(bytes)
	if err != nil {
		logger.Info().Msgf("Failed to read values: %v", err)
		return err
	}
	values := result.AsMap()
	pl := values["platform"].(map[string]interface{})
	plArtifacts := pl["artifacts"].(map[string]interface{})
	if platformSvcTag != "" {
		plArtifacts["tag"] = platformSvcTag
	}
	if platformSvcDigest != "" {
		plArtifacts["digest"] = platformSvcDigest
	}
	pl["artifacts"] = plArtifacts
	values["platform"] = pl

	vdas := values[VAPUSDATA_ARTIFACTS].(map[string]interface{})

	ctlArtifacts := vdas["vapusCtlArtifacts"].(map[string]interface{})
	if vapusctlSvcTag != "" {
		ctlArtifacts["tag"] = vapusctlSvcTag
	}
	if vapusctlSvcDigest != "" {
		ctlArtifacts["digest"] = vapusctlSvcDigest
	}
	vdas["vapusCtlArtifacts"] = ctlArtifacts

	operatorArtifacts := vdas["vapusOperatorArtifacts"].(map[string]interface{})
	if vapusoperatorSvcDigest != "" {
		operatorArtifacts["digest"] = vapusoperatorSvcDigest
	}
	if vapusoperatorSvcTag != "" {
		operatorArtifacts["tag"] = vapusoperatorSvcTag
	}
	vdas["vapusOperatorArtifacts"] = operatorArtifacts

	etlArtifacts := vdas["etlWorkerArtifacts"].(map[string]interface{})
	if etlworkerSvcDigest != "" {
		etlArtifacts["digest"] = etlworkerSvcDigest
	}
	if etlworkerSvcTag != "" {
		etlArtifacts["tag"] = etlworkerSvcTag
	}
	vdas["etlWorkerArtifacts"] = etlArtifacts

	vdcArtifacts := vdas["vapusDcArtifacts"].(map[string]interface{})
	if vapusDcSvcDigest != "" {
		vdcArtifacts["digest"] = vapusDcSvcDigest
	}
	if vapusDcSvcTag != "" {
		vdcArtifacts["tag"] = vapusDcSvcTag
	}
	vdas["vapusDcArtifacts"] = vdcArtifacts

	values[VAPUSDATA_ARTIFACTS] = vdas

	logger.Info().Msgf("Values --------: %v", vdas)
	bytes, err = yaml.Marshal(values)
	if err != nil {
		logger.Err(err).Msgf("Failed to marshal values: %v", err)
		return err
	}
	err = os.WriteFile(file, bytes, 0644)
	if err != nil {
		logger.Err(err).Msgf("Failed to write values: %v", err)
	}
	return err
}
