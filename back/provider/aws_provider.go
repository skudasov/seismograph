package provider

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

const (
	DefaultLoadGeneratorInstanceType = "t2.micro"
	DefaultVMRootUser                = "ec2-user"
)

type InfrastructureProviderAWS struct {
	client           *ec2.EC2
	session          *session.Session
	ClusterSpec      ClusterSpec
	RunningInstances map[string]*RunningInstance
}

type RunningInstance struct {
	Id              string
	Name            string
	KeyFileName     string
	PrivateKeyPem   string
	PublicDNSName   string
	PublicIPAddress string
}

type CreatedInstances struct {
	Instances []*CreatedInstance
}

type CreatedInstance struct {
	ID              string
	Region          string
	Name            string
	Image           string
	Type            string
	PrivateKeyPem   string
	PublicDNSName   string
	PublicIPAddress string
}

type ClusterSpec struct {
	Region    string
	Instances []InstanceSpec
}

type InstanceSpec struct {
	Region string `json:"region"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Type   string `json:"type"`
}

func stopInstance(svc *ec2.EC2, id string) error {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
		DryRun: aws.Bool(true),
	}
	result, err := svc.StopInstances(input)
	awsErr, ok := err.(awserr.Error)
	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = svc.StopInstances(input)
		if err != nil {
			return err
		} else {
			log.Infof("instance stopped: %s", result.StoppingInstances)
			return err
		}
	} else {
		return err
	}
}

func (m *InfrastructureProviderAWS) createInstances() (*CreatedInstances, error) {
	createdInstances := make([]*CreatedInstance, 0)
	for _, instanceSpec := range m.ClusterSpec.Instances {
		privateKey := createKeyPair(m.client, instanceSpec.Name)
		runResult, err := m.client.RunInstances(&ec2.RunInstancesInput{
			ImageId:      aws.String(instanceSpec.Image),
			InstanceType: aws.String(instanceSpec.Type),
			MinCount:     aws.Int64(1),
			MaxCount:     aws.Int64(1),
			KeyName:      aws.String(instanceSpec.Name),
			// must be security group with tcp 22, 8181 allowance, configured once for aws account
			// SecurityGroupIds: []*string{aws.String("sg-6dc1de0b")},
		})
		if err != nil {
			return &CreatedInstances{}, err
		}
		if len(runResult.Instances) == 0 || len(runResult.Instances) > 1 {
			return &CreatedInstances{}, errors.Wrap(err, "failed to start ec2 instance")
		}
		id := runResult.Instances[0].InstanceId
		log.Infof("instance created: %s (aws_id: %s)", instanceSpec.Name, *id)
		_, errtag := m.client.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{runResult.Instances[0].InstanceId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(instanceSpec.Name),
				},
				{
					Key:   aws.String("Id"),
					Value: id,
				},
			},
		})
		if errtag != nil {
			return &CreatedInstances{}, errors.New("could not create tags for instance")
		}
		log.Infof("instance tagged: %s", *id)
		createdInstances = append(createdInstances, &CreatedInstance{
			ID:            *id,
			Region:        instanceSpec.Region,
			Name:          instanceSpec.Name,
			Image:         instanceSpec.Image,
			Type:          instanceSpec.Type,
			PrivateKeyPem: privateKey,
			// must be acquired later, when vm in state "running"
			PublicDNSName:   "",
			PublicIPAddress: "",
		})
	}
	return &CreatedInstances{Instances: createdInstances}, nil
}

func (m *InfrastructureProviderAWS) collectPublicAddresses(instances *CreatedInstances) {
	for _, r := range instances.Instances {
		res := DescribeInstances(m.client, filterById(r.ID))
		r.PublicIPAddress = *res.Reservations[0].Instances[0].PublicIpAddress
		r.PublicDNSName = *res.Reservations[0].Instances[0].PublicDnsName
		log.Debugf("addresses assigned: %s (%s)", r.PublicDNSName, r.PublicIPAddress)
	}
}

func (m *InfrastructureProviderAWS) provision() {
	for _, r := range m.RunningInstances {
		userDnsString := fmt.Sprintf("%s@%s", DefaultVMRootUser, r.PublicDNSName)
		log.Infof("connection string: \nssh -i %s %s", userDnsString, userDnsString)
	}
}

func (m *InfrastructureProviderAWS) dumpPrivateKeysToFile(instances *CreatedInstances) {
	for _, r := range instances.Instances {
		dumpPrivateKeyPem(fmt.Sprintf("%s@%s", DefaultVMRootUser, r.PublicDNSName), r.PrivateKeyPem)
	}
}

// Bootstrap creates vms according to spec and wait until all vm in state "running"
func (m *InfrastructureProviderAWS) Bootstrap() *CreatedInstances {
	instances, err := m.createInstances()
	if err != nil {
		log.Error("failed to create instances: %s", err)
	}
	for i := 0; i < 50; i++ {
		time.Sleep(5 * time.Second)
		if m.assureRunning() {
			log.Info("all instances are running")
			break
		}
	}
	m.collectPublicAddresses(instances)
	m.dumpPrivateKeysToFile(instances)
	// TODO: general retry function, even if VM is in state "running" ssh may be unavailable,
	//  has no flag to know it for sure
	return instances
}

func filterById(id string) *ec2.DescribeInstancesInput {
	return &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Id"),
				Values: []*string{aws.String(id)},
			},
		},
	}
}

func (m *InfrastructureProviderAWS) assureRunning() bool {
	values := make([]*string, 0)
	for _, i := range m.RunningInstances {
		values = append(values, aws.String(i.Id))
	}
	filter := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Id"),
				Values: values,
			},
		},
	}
	res := DescribeInstances(m.client, filter)
	for _, r := range res.Reservations {
		if *r.Instances[0].State.Name != "running" {
			return false
		}
	}
	return true
}

func (m *InfrastructureProviderAWS) Exec(vmName string, cmd string) {
	for _, r := range m.RunningInstances {
		if r.Name == vmName {
			client, sess, err := connectSSHToHost(DefaultVMRootUser, r.PublicIPAddress, r.PrivateKeyPem)
			if err != nil {
				panic(err)
			}
			log.Debugf("executing cmd on vm %s: %s", r.Id, cmd)
			out, err := sess.CombinedOutput(cmd)
			if err != nil {
				log.Fatal(err)
			}
			log.Debugf(string(out))
			client.Close()
		}
	}
}

func newEC2Session(region string) (*ec2.EC2, *session.Session) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:                        aws.String(region),
			CredentialsChainVerboseErrors: aws.Bool(true),
			// for now creds mounted in compose from host ~/.aws/config and ~/.aws/credentials
			// Credentials: credentials.NewStaticCredentials(),
		},
		SharedConfigState: session.SharedConfigEnable,
	},
	)
	if err != nil {
		log.Fatal(err)
	}
	return ec2.New(sess), sess
}

func createKeyPair(svc *ec2.EC2, name string) string {
	result, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(name),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Created key pair %q %s\n%s\n",
		*result.KeyName, *result.KeyFingerprint,
		*result.KeyMaterial)
	return *result.KeyMaterial
}

func NewInfrastructureProviderAWS(spec ClusterSpec) *InfrastructureProviderAWS {
	client, sess := newEC2Session(spec.Region)
	return &InfrastructureProviderAWS{
		session:          sess,
		client:           client,
		ClusterSpec:      spec,
		RunningInstances: make(map[string]*RunningInstance, 0),
	}
}

func connectSSHToHost(user, ip, privateKeyPem string) (*ssh.Client, *ssh.Session, error) {
	signer, err := ssh.ParsePrivateKey([]byte(privateKeyPem))
	if err != nil {
		log.Fatal(err)
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", ip+":22", sshConfig)
	if err != nil {
		return nil, nil, err
	}

	sess, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, sess, nil
}

func dumpPrivateKeyPem(name string, key string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(key)); err != nil {
		log.Fatal(err)
	}
	// ssh requirements to use
	if err := os.Chmod(name, 0400); err != nil {
		log.Fatal(err)
	}
}

func DescribeInstances(svc *ec2.EC2, input *ec2.DescribeInstancesInput) *ec2.DescribeInstancesOutput {
	result, err := svc.DescribeInstances(input)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
