package server

import (
	"github.com/f4hrenh9it/seismograph/back/provider"
)

type ProjectDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RepoUrl     string `json:"repo_url"`
}

type GetProjectsResponse struct {
	ErrorResponse
	Projects []ProjectDTO `json:"projects"`
}

type AttackClusterRequest struct {
	Name         string               `json:"name"`
	ProviderName string               `json:"provider_name"`
	ProjectID    uint                 `json:"project_id"`
	ClusterSpec  provider.ClusterSpec `json:"cluster_spec"`
}

type AttackClusterResponse struct {
	ErrorResponse
	ID uint `json:"id"`
}

type GetAttackClusterResponse struct {
	ErrorResponse
	AttackClusterDTO
}

type GetAttackClusterVMSResponse struct {
	ErrorResponse
	Instances []Instance `json:"instances"`
}

type GetAttackClustersResponse struct {
	ErrorResponse
	AttackClustersDTO []AttackClusterDTO `json:"clusters"`
}

type GetAttackVMResponse struct {
	ErrorResponse
	AttackVMDTO
}

type AttackVMDTO struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	AttackClusterID uint   `json:"attack_cluster_id"`
}

type AttackVMRequest struct {
	Name            string `json:"name"`
	AttackClusterID uint   `json:"attack_cluster_id"`
}

type AttackVMResponse struct {
	ErrorResponse
	ID uint `json:"id"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	RepoURL     string `json:"repo_url"`
}

type CreateProjectResponse struct {
	ErrorResponse
	ID uint `json:"id"`
}

type AddTestDataResponse struct {
	ErrorResponse
	DataID uint   `json:"data_id"`
	BlobID string `json:"blob_id"`
}

type TestDataDTO struct {
	ID           uint   `json:"id"`
	TestName     string `json:"test_name"`
	Environment  string `json:"environment"`
	ProjectID    uint   `json:"project_id"`
	DataBlobPath string `json:"data_blob_path"`
}

type AttackClusterDTO struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	ProjectID    uint   `json:"project_id"`
}

type TestsMetaResponse struct {
	ErrorResponse
	Tests []TestDataDTO `json:"tests"`
}

type TestDataResponse struct {
	ErrorResponse
	Chart TestChart `json:"chart"`
}

type Series map[string][]string

type TestChart struct {
	TestName     string `json:"test_name"`
	Environment  string `json:"environment"`
	DataBlobPath string `json:"data_blob_path"`
	Series       Series `json:"series"`
}
