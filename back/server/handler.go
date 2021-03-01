package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/f4hrenh9it/seismograph/back/provider"
	"github.com/f4hrenh9it/seismograph/back/types"
)

func (m *SeismographService) HandleCreateProject(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	var projBody CreateProjectRequest
	if err := c.BodyParser(&projBody); err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	projID, err := m.StoreProject(ctx, projBody)
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	return c.JSON(&CreateProjectResponse{
		ID: projID,
	})
}

func (m *SeismographService) HandleProjects(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	projects, err := m.GetProjects(ctx)
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	ps := make([]ProjectDTO, 0)
	for _, p := range projects {
		ps = append(ps, p.DTO())
	}
	return c.JSON(&GetProjectsResponse{
		Projects: ps,
	})
}

func (m *SeismographService) HandleDeleteProjectWithTests(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	projId := c.Params("project")
	pid, _ := strconv.Atoi(projId)
	err := m.DeleteProjectWithTests(ctx, uint(pid))
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	return c.SendStatus(200)
}

func (m *SeismographService) HandleTestsForProject(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	projId := c.Params("project")
	pid, _ := strconv.Atoi(projId)
	tests, err := m.LoadTestsDataForProject(ctx, uint(pid))
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	ts := make([]TestDataDTO, 0)
	for _, t := range tests {
		ts = append(ts, t.DTO())
	}
	return c.JSON(&TestsMetaResponse{
		Tests: ts,
	})
}

func (m *SeismographService) HandleGetTestDataId(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	id := c.Params("id")
	projId := c.Params("project")
	pid, _ := strconv.Atoi(projId)
	chart, err := m.LoadTestChartByDataID(ctx, uint(pid), id)
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	return c.JSON(&TestDataResponse{
		Chart: *chart,
	})
}

func (m *SeismographService) HandleDeleteTest(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	id := c.Params("id")
	projId := c.Params("project")
	pid, _ := strconv.Atoi(projId)
	if err := m.DeleteTest(ctx, uint(pid), id); err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	return c.SendStatus(200)
}

func (m *SeismographService) HandleTest(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()
	l := log.WithContext(ctx)

	mpf, err := c.MultipartForm()
	if err != nil {
		log.Fatal(err)
	}
	projId := c.Params("project")
	pid, _ := strconv.Atoi(projId)
	l.WithField("meta", mpf.Value)
	if err := ValidateAddDataRequest(mpf); err != "" {
		return c.JSON(DefaultErrorResponse(err))
	}
	var meta map[string]string
	if err := json.Unmarshal([]byte(mpf.Value["meta"][0]), &meta); err != nil {
		return c.JSON(DefaultErrorResponse(MetaParsingErr))
	}
	if err := m.CreateBucket(ctx, types.DefaultBucketName); err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	blobPath := fmt.Sprintf("test-data-%s", uuid.New().String())
	for _, fh := range mpf.File[types.MPartDataKey] {
		f, err := fh.Open()
		if err != nil {
			return c.JSON(DefaultErrorResponse(MultipartReadErr))
		}
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, f); err != nil {
			log.Fatal(err)
		}
		l.WithFields(log.Fields{"blob_path": blobPath}).Infof("storing blob")
		if err := m.StoreBlob(ctx, types.DefaultBucketName, blobPath, fh.Size, buf); err != nil {
			return c.JSON(DefaultErrorResponse(TestDataCreationErr))
		}
	}
	l.Infof("storing meta info for blob: %s", blobPath)
	dataID, blobID, err := m.StoreTestData(ctx, uint(pid), meta, blobPath)
	if err != nil {
		return c.JSON(DefaultErrorResponse(MetaInfoCreationErr))
	}
	return c.JSON(&AddTestDataResponse{
		DataID: dataID,
		BlobID: blobID,
	})
}

func (m *SeismographService) HandleAttackCluster(c *fiber.Ctx) error {

	var acBody AttackClusterRequest
	if err := c.BodyParser(&acBody); err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	if err := ValidateClusterSpec(acBody); err != "" {
		return c.JSON(DefaultErrorResponse(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Cluster.CreationTimeout)
	g, ctx := errgroup.WithContext(ctx)
	acID, err := m.StoreAttackCluster(ctx, acBody)
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}

	g.Go(func() error {
		prov := provider.NewInfrastructureProviderAWS(acBody.ClusterSpec)
		inst := prov.Bootstrap()

		if err := m.UpdateAttackClusterAfterCreation(ctx, acID, acBody, *inst); err != nil {
			return err
		}
		cancel()
		return nil
	})
	return c.JSON(&AttackClusterResponse{
		ID: acID,
	})
}

func (m *SeismographService) HandleGetAttackCluster(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	id := c.Params("id")
	cid, _ := strconv.Atoi(id)
	ac, err := m.GetAttackCluster(ctx, uint(cid))
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	return c.JSON(&GetAttackClusterResponse{
		AttackClusterDTO: ac.DTO(),
	})
}

func (m *SeismographService) HandleGetAttackClusterVMS(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	id := c.Params("id")
	cid, _ := strconv.Atoi(id)
	vms, err := m.GetAttackClusterVMS(ctx, uint(cid))
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	return c.JSON(&GetAttackClusterVMSResponse{
		Instances: vms,
	})
}

func (m *SeismographService) HandleGetAttackClusters(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	clusters, err := m.GetAttackClusters(ctx)
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	acs := make([]AttackClusterDTO, 0)
	for _, ac := range clusters {
		acs = append(acs, ac.DTO())
	}
	return c.JSON(&GetAttackClustersResponse{
		AttackClustersDTO: acs,
	})
}

func (m *SeismographService) HandleAttackVM(c *fiber.Ctx) error {
	var vmBody AttackVMRequest
	if err := c.BodyParser(&vmBody); err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()
	id, err := m.StoreAttackVM(ctx, vmBody)
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	return c.JSON(&AttackVMResponse{
		ID: id,
	})
}

func (m *SeismographService) HandleGetAttackVM(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	id := c.Params("id")
	cid, _ := strconv.Atoi(id)
	vm, err := m.GetAttackVM(ctx, uint(cid))
	if err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	return c.JSON(&GetAttackVMResponse{
		AttackVMDTO: vm.DTO(),
	})
}

func (m *SeismographService) HandleDeleteAttackVM(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	id := c.Params("id")
	vmid, _ := strconv.Atoi(id)
	if err := m.DeleteAttackVM(ctx, uint(vmid)); err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	return c.SendStatus(200)
}

func (m *SeismographService) HandleDeleteAttackCluster(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.Server.RequestTimeout)
	defer cancel()

	id := c.Params("id")
	cid, _ := strconv.Atoi(id)
	if err := m.DeleteAttackCluster(ctx, uint(cid)); err != nil {
		return c.JSON(DefaultErrorResponse(err.Error()))
	}
	return c.SendStatus(200)
}
