# Makefile

# Go parameters
GOCMD = go
BUFCMD = buf

.PHONY: protogen gosync api

protogen:
	$(BUFCMD) lint
	$(BUFCMD) generate
	protoc-go-inject-tag -input gen-go/protos/aistudio/v1alpha1/*.pb.go
	protoc-go-inject-tag -input gen-go/protos/models/v1alpha1/common.pb.go


gosync: protogen
	cd gen-go && $(GOCMD) mod tidy && $(GOCMD) mod vendor

api: gosync