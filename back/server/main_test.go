package server

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/f4hrenh9it/seismograph/back/config"
	"github.com/f4hrenh9it/seismograph/back/migrations"
	"github.com/f4hrenh9it/seismograph/back/provider"
	"github.com/f4hrenh9it/seismograph/back/testutil"
	"github.com/f4hrenh9it/seismograph/back/types"
)

const (
	baseURL = "http://localhost:10500"
)

var (
	DefaultMetaInfoBinary []byte
)

func prepareDefaultTestData() {
	DefaultMetaInfoBinary, _ = json.Marshal(map[string]interface{}{
		"test_name":   "test_1",
		"environment": "load1",
	})
}

func createProject(t *testing.T) CreateProjectResponse {
	e := httpexpect.New(t, baseURL)
	body := e.POST(URLHandleCreateProject).
		WithJSON(CreateProjectRequest{Name: "test_proj_1"}).
		Expect().
		Status(200).
		Body().
		Raw()
	var r CreateProjectResponse
	err := json.Unmarshal([]byte(body), &r)
	require.NoError(t, err)
	return r
}

func createAttackCluster(t *testing.T, pid uint, req *AttackClusterRequest) AttackClusterResponse {
	e := httpexpect.New(t, baseURL)
	var body AttackClusterRequest
	if req == nil {
		body = AttackClusterRequest{
			Name:         "test_cluster",
			ProviderName: "aws",
			ProjectID:    pid,
			ClusterSpec: provider.ClusterSpec{
				Region: CFG.Cluster.DefaultRegion,
				Instances: []provider.InstanceSpec{
					{
						Name:   uuid.New().String(),
						Region: CFG.Cluster.DefaultRegion,
						Image:  CFG.Cluster.DefaultImage,
						Type:   provider.DefaultLoadGeneratorInstanceType,
					},
					{
						Name:   uuid.New().String(),
						Region: CFG.Cluster.DefaultRegion,
						Image:  CFG.Cluster.DefaultImage,
						Type:   provider.DefaultLoadGeneratorInstanceType,
					},
				},
			},
		}
	} else {
		body = *req
	}
	b := e.POST(URLHandleCreateAttackCluster).
		WithJSON(body).
		Expect().
		Status(200).
		Body().
		Raw()
	var r AttackClusterResponse
	err := json.Unmarshal([]byte(b), &r)
	require.NoError(t, err)
	return r
}

func createAttackVM(t *testing.T, cid uint) AttackVMResponse {
	e := httpexpect.New(t, baseURL)
	body := AttackVMRequest{
		Name:            "test_vm_1",
		AttackClusterID: cid,
	}
	b := e.POST(URLHandleCreateAttackVM).
		WithJSON(body).
		Expect().
		Status(200).
		Body().
		Raw()
	var r AttackVMResponse
	err := json.Unmarshal([]byte(b), &r)
	require.NoError(t, err)
	return r
}

func createTestInstance(t *testing.T, pid string) AddTestDataResponse {
	e := httpexpect.New(t, baseURL)
	body := e.POST(baseURL+fmt.Sprintf("/project/%s/test", pid)).
		WithMultipart().
		WithFile(types.MPartDataKey, "../testdata/testdata.csv").
		WithFormField("meta", string(DefaultMetaInfoBinary)).
		Expect().
		Status(200).
		Body().
		Raw()
	var resObj AddTestDataResponse
	err := json.Unmarshal([]byte(body), &resObj)
	require.NoError(t, err)
	return resObj
}

var CFG *config.Config

func TestMain(m *testing.M) {
	pool, resource, pgPort := testutil.RunPGSQL()
	poolM, resourceM, minioPort := testutil.RunMinio()
	CFG = &config.Config{
		Cluster: config.ClusterConfig{
			CreationTimeout: 2 * time.Minute,
			DefaultRegion:   "us-east-2",
			DefaultImage:    "ami-0c0415cdff14e2a4a",
		},
		Server: config.ServerConfig{
			Port:           "10500",
			RequestTimeout: 3 * time.Second,
		},
		DB: config.DBConfig{
			Postgres: config.PostgresConfig{
				Host:   "localhost",
				User:   "postgres",
				Pwd:    "secret",
				DBName: "testdb",
				Port:   pgPort,
			},
			Minio: config.MinioConfig{
				Url:       "localhost:9000",
				AccessKey: "minioadmin",
				SecretKey: "minioadmin",
				Port:      minioPort,
				Path:      "/data/minio",
			},
		},
	}
	go NewSeismographService(CFG)
	migrations.GormMigrateInit(CFG.DB)
	time.Sleep(1 * time.Second)
	prepareDefaultTestData()
	code := m.Run()
	_ = pool.Purge(resource)
	_ = poolM.Purge(resourceM)
	os.Exit(code)
}
