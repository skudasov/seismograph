package server

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	RepoUrl     string
}

func (p *Project) DTO() ProjectDTO {
	return ProjectDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		RepoUrl:     p.RepoUrl,
	}
}

type TestData struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	TestName     string
	Environment  string
	ProjectID    uint `gorm:"foreignKey:ID;references:Project"`
	DataBlobPath string
}

func (p *TestData) DTO() TestDataDTO {
	return TestDataDTO{
		ID:           p.ID,
		TestName:     p.TestName,
		Environment:  p.Environment,
		ProjectID:    p.ProjectID,
		DataBlobPath: p.DataBlobPath,
	}
}

type Instance struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"`
	ClusterID     uint `gorm:"foreignKey:ID;references:AttackCluster"`
	Region        string
	Name          string
	PublicDNSName string
	PrivateKeyPEM string
	Image         string
	Type          string
}

type AttackCluster struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	Name         string
	Region       string
	ProviderName string
	ProjectID    uint `gorm:"foreignKey:ID;references:Project"`
}

func (p *AttackCluster) DTO() AttackClusterDTO {
	return AttackClusterDTO{
		ID:           p.ID,
		Name:         p.Name,
		ProviderName: p.ProviderName,
		ProjectID:    p.ProjectID,
	}
}

type AttackVM struct {
	gorm.Model
	ID              uint `gorm:"primaryKey"`
	Name            string
	AttackClusterID uint `gorm:"foreignKey:ID;references:AttackCluster"`
}

func (p *AttackVM) DTO() AttackVMDTO {
	return AttackVMDTO{
		ID:              p.ID,
		Name:            p.Name,
		AttackClusterID: p.AttackClusterID,
	}
}
