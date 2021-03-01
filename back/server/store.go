package server

import (
	"bytes"
	"context"
	"io"
	"strconv"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"

	"github.com/f4hrenh9it/seismograph/back/provider"
	"github.com/f4hrenh9it/seismograph/back/types"
)

func (m *SeismographService) StoreAttackVM(ctx context.Context, req AttackVMRequest) (uint, error) {
	vm := AttackVM{
		Name:            req.Name,
		AttackClusterID: req.AttackClusterID,
	}
	if err := m.PG.WithContext(ctx).Create(&vm).Error; err != nil {
		return 0, err
	}
	return vm.ID, nil
}

func (m *SeismographService) StoreAttackCluster(ctx context.Context, req AttackClusterRequest) (uint, error) {
	ac := AttackCluster{
		Name:         req.Name,
		ProviderName: req.ProviderName,
		ProjectID:    req.ProjectID,
		Region:       req.ClusterSpec.Region,
	}
	if err := m.PG.WithContext(ctx).Create(&ac).Error; err != nil {
		return 0, err
	}
	return ac.ID, nil
}

func (m *SeismographService) UpdateAttackClusterAfterCreation(ctx context.Context, id uint, req AttackClusterRequest, inst provider.CreatedInstances) error {
	if err := m.PG.Transaction(func(tx *gorm.DB) error {
		for _, i := range inst.Instances {
			instanceModel := Instance{
				ClusterID:     id,
				Region:        req.ClusterSpec.Region,
				Name:          i.Name,
				PublicDNSName: i.PublicDNSName,
				PrivateKeyPEM: i.PrivateKeyPem,
				Image:         i.Image,
				Type:          i.Type,
			}
			if err := tx.WithContext(ctx).Create(&instanceModel).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (m *SeismographService) GetAttackCluster(ctx context.Context, id uint) (AttackCluster, error) {
	ac := AttackCluster{
		ID: id,
	}
	if err := m.PG.WithContext(ctx).Find(&ac).Error; err != nil {
		return AttackCluster{}, err
	}
	return ac, nil
}

func (m *SeismographService) GetAttackClusterVMS(ctx context.Context, id uint) ([]Instance, error) {
	var instances []Instance
	if err := m.PG.WithContext(ctx).Where("cluster_id = ?", id).Find(&instances).Error; err != nil {
		return []Instance{}, err
	}
	return instances, nil
}

func (m *SeismographService) GetAttackClusters(ctx context.Context) ([]AttackCluster, error) {
	var acs []AttackCluster
	if err := m.PG.WithContext(ctx).Find(&acs).Error; err != nil {
		return []AttackCluster{}, err
	}
	return acs, nil
}

func (m *SeismographService) DeleteAttackCluster(ctx context.Context, id uint) error {
	data := AttackCluster{
		ID: id,
	}
	if err := m.PG.WithContext(ctx).Where("id = ?", id).Delete(&data).Error; err != nil {
		return err
	}
	return nil
}

func (m *SeismographService) GetAttackVM(ctx context.Context, id uint) (AttackVM, error) {
	vm := AttackVM{
		ID: id,
	}
	if err := m.PG.WithContext(ctx).Find(&vm).Error; err != nil {
		return AttackVM{}, err
	}
	return vm, nil
}

func (m *SeismographService) GetProjects(ctx context.Context) ([]Project, error) {
	var projs []Project
	if err := m.PG.WithContext(ctx).Find(&projs).Error; err != nil {
		return projs, err
	}
	return projs, nil
}

func (m *SeismographService) StoreProject(ctx context.Context, req CreateProjectRequest) (uint, error) {
	proj := Project{
		Name:        req.Name,
		Description: req.Description,
		RepoUrl:     req.RepoURL,
	}
	if err := m.PG.WithContext(ctx).Create(&proj).Error; err != nil {
		return 0, err
	}
	return proj.ID, nil
}

func (m *SeismographService) CreateBucket(ctx context.Context, bucketName string) error {
	found, err := m.M.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !found {
		if err := m.M.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (m *SeismographService) StoreBlob(ctx context.Context, bucketName, objName string, objSize int64, d io.Reader) error {
	_, err := m.M.PutObject(ctx, bucketName, objName, d, objSize, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *SeismographService) GetBlob(ctx context.Context, bucketName, objName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	o, err := m.M.GetObject(ctx, bucketName, objName, opts)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (m *SeismographService) LoadTestsDataForProject(ctx context.Context, pid uint) ([]TestData, error) {
	var td []TestData
	if err := m.PG.WithContext(ctx).Where("project_id = ?", pid).Find(&td).Error; err != nil {
		return nil, err
	}
	return td, nil
}

func (m *SeismographService) DeleteTest(ctx context.Context, pid uint, tid string) error {
	tID, _ := strconv.Atoi(tid)
	data := TestData{
		ID: uint(tID),
	}
	if err := m.PG.WithContext(ctx).Where("project_id = ?", pid).Delete(&data).Error; err != nil {
		return err
	}
	return nil
}

func (m *SeismographService) LoadTestChartByDataID(ctx context.Context, pid uint, tid string) (*TestChart, error) {
	tID, _ := strconv.Atoi(tid)
	data := TestData{
		ID: uint(tID),
	}
	if err := m.PG.WithContext(ctx).Where("project_id = ?", pid).Find(&data).Error; err != nil {
		return nil, err
	}
	csv, err := m.GetBlob(ctx, types.DefaultBucketName, data.DataBlobPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, csv); err != nil {
		return nil, err
	}
	return &TestChart{
		TestName:     data.TestName,
		Environment:  data.Environment,
		DataBlobPath: data.DataBlobPath,
		Series:       DataToKVSlices(buf),
	}, nil
}

func (m *SeismographService) StoreTestData(ctx context.Context, projID uint, meta map[string]string, blobPath string) (uint, string, error) {
	td := &TestData{
		TestName:     meta["test_name"],
		Environment:  meta["environment"],
		ProjectID:    projID,
		DataBlobPath: blobPath,
	}
	if err := m.PG.WithContext(ctx).Create(&td).Error; err != nil {
		return 0, "", err
	}
	return td.ID, blobPath, nil
}

func (m *SeismographService) DeleteProjectWithTests(ctx context.Context, pid uint) error {
	if err := m.PG.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Where("project_id = ?", pid).Delete(&TestData{}).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Where("id = ?", pid).Delete(&Project{}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (m *SeismographService) DeleteAttackVM(ctx context.Context, vmid uint) error {
	data := AttackVM{
		ID: vmid,
	}
	if err := m.PG.WithContext(ctx).Where("id = ?", vmid).Delete(&data).Error; err != nil {
		return err
	}
	return nil
}
