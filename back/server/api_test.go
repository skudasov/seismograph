package server

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/google/uuid"

	"github.com/f4hrenh9it/seismograph/back/provider"
	"github.com/f4hrenh9it/seismograph/back/types"
)

func TestDataUpload(t *testing.T) {
	p := createProject(t)
	pid := strconv.FormatUint(uint64(p.ID), 10)
	e := httpexpect.New(t, baseURL)

	t.Run("create upload", func(t *testing.T) {
		res := e.POST(fmt.Sprintf("/project/%s/test", pid)).
			WithMultipart().
			WithFile(types.MPartDataKey, "../testdata/testdata.csv").
			WithFormField("meta", string(DefaultMetaInfoBinary)).
			Expect().
			JSON().
			Object()
		res.Value("data_id").NotNull()
		res.Value("blob_id").NotNull()
	})
	t.Run("no upload csv data", func(t *testing.T) {
		res := e.POST(fmt.Sprintf("/project/%s/test", pid)).
			WithMultipart().
			WithFile(types.MPartDataKey, "../testdata/empty.csv").
			WithFormField("meta", string(DefaultMetaInfoBinary)).
			Expect().
			JSON().
			Object()
		res.Value("errors").Equal([]string{EmptyBlobErr})
	})

	t.Run("meta parsing failed", func(t *testing.T) {
		e.POST(fmt.Sprintf("/project/%s/test", pid)).
			WithMultipart().
			WithFile(types.MPartDataKey, "../testdata/testdata.csv").
			WithFormField("meta", []byte("")).
			Expect().
			JSON().
			Object().Value("errors").Equal([]string{MetaParsingErr})
	})

	t.Run("no meta field", func(t *testing.T) {
		e.POST(fmt.Sprintf("/project/%s/test", pid)).
			WithMultipart().
			WithFile(types.MPartDataKey, "../testdata/testdata.csv").
			Expect().
			JSON().
			Object().Value("errors").Equal([]string{NoMetaFieldErr})
	})
}

func TestTestDataCRUD(t *testing.T) {
	e := httpexpect.New(t, baseURL)
	p := createProject(t)
	pid := strconv.FormatUint(uint64(p.ID), 10)
	t.Run("get projects", func(t *testing.T) {
		e.GET(URLHandleProjects).Expect().JSON().Object().Value("projects").NotNull()
	})
	t.Run("get tests", func(t *testing.T) {
		createTestInstance(t, pid)
		e.GET(fmt.Sprintf("/project/%s/tests", pid)).
			Expect().JSON().Object().Value("tests").Array().NotEmpty()
	})
	t.Run("get test data by id", func(t *testing.T) {
		createdTest := createTestInstance(t, pid)
		e.GET(fmt.Sprintf("/project/%s/test/%d", pid, createdTest.DataID)).
			Expect().
			JSON().
			Object().
			Value("chart").
			Equal(
				TestChart{
					TestName:     "test_1",
					Environment:  "load1",
					DataBlobPath: createdTest.BlobID,
					Series: Series{
						"BeginTimeNano": {"1601972898498711000", "1601972898625900000"},
						"EndTimeNano":   {"1601972899003622000", "1601972899126288000"},
						"Error":         {"ok", "ok"},
						"RequestLabel":  {"test_runner", "test_runner"},
					},
				})
	})
}

func TestDeleteProject(t *testing.T) {
	e := httpexpect.New(t, baseURL)
	p := createProject(t)
	pid := strconv.FormatUint(uint64(p.ID), 10)
	createTestInstance(t, pid)
	e.DELETE(fmt.Sprintf("/project/%s", pid)).
		Expect().
		Status(200)
}

func TestDeleteSingleTest(t *testing.T) {
	e := httpexpect.New(t, baseURL)
	p := createProject(t)
	pid := strconv.FormatUint(uint64(p.ID), 10)
	i := createTestInstance(t, pid)
	e.DELETE(fmt.Sprintf("/project/%s/test/%d", pid, i.DataID)).
		Expect().
		Status(200)
}

func TestAttackClusterValidation(t *testing.T) {
	e := httpexpect.New(t, baseURL)
	p := createProject(t)
	t.Run("provider is not supported", func(t *testing.T) {
		e.POST(URLHandleCreateAttackCluster).
			WithJSON(&AttackClusterRequest{
				Name:         "test_cluster",
				ProviderName: "nothing",
				ProjectID:    p.ID,
				ClusterSpec: provider.ClusterSpec{
					Region: CFG.Cluster.DefaultRegion,
					Instances: []provider.InstanceSpec{
						{
							Name:   uuid.New().String(),
							Region: CFG.Cluster.DefaultRegion,
							Image:  CFG.Cluster.DefaultImage,
							Type:   provider.DefaultLoadGeneratorInstanceType,
						},
					},
				},
			}).
			Expect().
			JSON().Object().Value("errors").Equal([]string{UnsupportedProviderErr})
	})

	t.Run("bad image", func(t *testing.T) {
		e.POST(URLHandleCreateAttackCluster).
			WithJSON(&AttackClusterRequest{
				Name:         "test_cluster",
				ProviderName: "nothing",
				ProjectID:    p.ID,
				ClusterSpec: provider.ClusterSpec{
					Region: CFG.Cluster.DefaultRegion,
					Instances: []provider.InstanceSpec{
						{
							Name:   uuid.New().String(),
							Region: CFG.Cluster.DefaultRegion,
							Image:  "bad image",
							Type:   provider.DefaultLoadGeneratorInstanceType,
						},
					},
				},
			}).
			Expect().
			JSON().Object().Value("errors").Equal([]string{UnsupportedProviderErr})
	})
}

func TestAttackClusterCRUD(t *testing.T) {
	e := httpexpect.New(t, baseURL)
	p := createProject(t)
	c := createAttackCluster(t, p.ID, nil)
	e.GET(fmt.Sprintf("/attack_cluster/%d", c.ID)).
		Expect().
		JSON().
		Object().
		Value("id").
		Equal(1)
	e.GET(URLHandleGetAttackClusters).
		Expect().
		JSON().
		Object().
		Value("clusters").
		Array().
		First().
		Object().
		Value("id").
		Equal(1)
	time.Sleep(30 * time.Second)
	e.GET(fmt.Sprintf("/attack_cluster/%d/vms", c.ID)).
		Expect().
		JSON().
		Object().
		Value("instances").
		Array().
		NotEmpty()
	e.DELETE(fmt.Sprintf("/attack_cluster/%d", c.ID)).
		Expect().
		Status(200)
}

func TestAttackVMCRUD(t *testing.T) {
	p := createProject(t)
	c := createAttackCluster(t, p.ID, nil)
	vm := createAttackVM(t, c.ID)
	e := httpexpect.New(t, baseURL)
	e.GET(fmt.Sprintf("/vm/%d", vm.ID)).
		Expect().
		JSON().
		Object().
		Value("id").
		Equal(1)
	e.DELETE(fmt.Sprintf("/vm/%d", vm.ID)).
		Expect().
		Status(200)
}
