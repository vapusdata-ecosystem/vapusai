package gcp

import "errors"

var (
	ErrCreatingGcpArClient      = errors.New("error creating gcp artifact registry client")
	ErrCreatingGcpSMClient      = errors.New("error while creating GCP Secret Manager client")
	ErrReadingGcpSecret         = errors.New("error while reading GCP Secret Manager")
	ErrDeletingGcpSecret        = errors.New("error while deleting GCP Secret Manager")
	ErrCreatingGcpSecret        = errors.New("error while creating secret in GCP Secret Manager")
	ErrParsingGAR               = errors.New("error while parsing GCP Artifact Registry URL")
	ErrParsingGARRegion         = errors.New("error while parsing GCP Artifact Registry Region")
	ErrParsingGARHost           = errors.New("error while parsing GCP Artifact Registry Host")
	ErrListingGARPackages       = errors.New("error while listing GCP Artifact Registry packages")
	ErrCreatingGcpSecretVersion = errors.New("error while creating secret version in GCP Secret Manager")
)
