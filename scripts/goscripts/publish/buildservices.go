package vapuspublish

import (
	"bytes"
	"os"
	"os/exec"
)

// func ReleasePlatformSvc(opts *cranetp.CraneClient) {
func ReleasePlatformSvc() {
	var err error
	orignalDir, err := os.Getwd()
	if err != nil {
		logger.Error().Err(err).Msg("Error getting current working directory")
		return
	}
	logger.Info().Msgf("Current working directory is %s", orignalDir)
	defer func(path string) {
		os.Chdir(path)
	}(orignalDir)

	err = os.Chdir(PLATFOR_SVC_PATH)
	if err != nil {
		logger.Error().Err(err).Msgf("Error changing directory to %s", PLATFOR_SVC_PATH)
		return
	}
	logger.Info().Msgf("Changed directory to %s", PLATFOR_SVC_PATH)
	cmd := exec.Command("make", "release")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		logger.Error().Err(err).Msg("Error building binary")
		return
	}
	logger.Info().Msg("Successfully released platform service")
	logger.Info().Msgf("Output: %s", out.String())
	// // dockerfilePath := filepath.Join(PLATFOR_SVC_PATH, "Dockerfile")
	// tempDir := os.TempDir()
	// if err != nil {
	// 	logger.Error().Err(err).Msg("Error creating temp directory")
	// }
	// logger.Info().Msgf("Temp directory created at %s", tempDir)

	// defer os.RemoveAll(tempDir)

	// imageTarPath := filepath.Join(tempDir, PLATFORM_TAR)
	// cmd = exec.Command("docker", "build", "-t", opts.GetFullOciURL(), ".")
	// cmd.Stdout = &out
	// cmd.Stderr = &out
	// if err := cmd.Run(); err != nil {
	// 	logger.Error().Err(err).Msgf("Error building docker image: %s", out.String())
	// }

	// cmd = exec.Command("docker", "save", "-o", imageTarPath, opts.GetFullOciURL())
	// cmd.Stdout = &out
	// cmd.Stderr = &out
	// if err := cmd.Run(); err != nil {
	// 	logger.Error().Err(err).Msgf("Error saving docker image to tarball: %s", out.String())
	// 	return
	// }

	// // Load the tarball as an OCI image
	// img, err := tarball.ImageFromPath(imageTarPath, nil)
	// if err != nil {
	// 	logger.Error().Err(err).Msgf("Error loading image from tarball %s", imageTarPath)
	// 	return
	// }

	// // Push the image to the registry
	// digest, err := opts.PushImage(img)
	// if err != nil {
	// 	logger.Error().Err(err).Msgf("Error pushing image %s", opts.GetFullOciURL())
	// 	return
	// }
}
