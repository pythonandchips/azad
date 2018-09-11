package awsinventory

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pythonandchips/azad/plugin"
)

// Name of plugin as used in playbook files
const Name = "aws"

// GetSchema returns the schema for the current plugin
//
// Inventory
// 	 ec2:
// 		 access_key_id (required: false): access key id for the aws account
//     secret_key (required: false): secret key for the aws account
//     connect_in (required: false): the property to use as the host name when connection to server.
// 		  														  This can be "PublicIpAddress" (default), "PrivateDnsName", "PublicDnsName" or "PrivateIpAddress"
//     region (required: false): region to search for resources. defaults to us-east-1
func GetSchema() plugin.Schema {
	return plugin.Schema{
		Inventory: map[string]plugin.Inventory{
			"ec2": {
				Fields: []plugin.Field{
					{Name: "access_key_id", Type: "String", Required: false},
					{Name: "secret_key", Type: "String", Required: false},
					{Name: "connect_on", Type: "String", Required: false},
					{Name: "region", Type: "String", Required: false},
				},
				Run: ec2Resources,
			},
		},
	}
}

const defaultRegion = "us-east-1"

type ec2iface interface {
	DescribeInstancesPages(*ec2.DescribeInstancesInput, func(*ec2.DescribeInstancesOutput, bool) bool) error
}

var ec2Session = func(context plugin.InventoryContext) ec2iface {
	sess := session.Must(
		session.NewSession(&aws.Config{
			MaxRetries: aws.Int(3),
		}),
	)
	region := context.GetWithDefault("region", defaultRegion)
	svc := ec2.New(sess, &aws.Config{
		Region: aws.String(region),
	})
	return svc
}

func ec2Resources(context plugin.InventoryContext) ([]plugin.Resource, error) {
	svc := ec2Session(context)
	filters := []*ec2.Filter{
		{
			Name:   aws.String("instance-state-name"),
			Values: aws.StringSlice([]string{"running"}),
		},
	}
	resources := []plugin.Resource{}
	request := ec2.DescribeInstancesInput{Filters: filters}
	err := svc.DescribeInstancesPages(
		&request,
		func(result *ec2.DescribeInstancesOutput, lastPage bool) bool {
			resources = parseResource(context, resources, result, lastPage)
			return lastPage
		},
	)
	return resources, err
}

func parseResource(
	context plugin.InventoryContext,
	resources []plugin.Resource,
	result *ec2.DescribeInstancesOutput,
	lastPage bool,
) []plugin.Resource {
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			resource := plugin.Resource{}
			connectOn := context.GetWithDefault("connect_on", "PublicIpAddress")
			switch connectOn {
			case "PublicIpAddress":
				resource.ConnectOn = *instance.PublicIpAddress
			case "PrivateDnsName":
				resource.ConnectOn = *instance.PrivateDnsName
			case "PublicDnsName":
				resource.ConnectOn = *instance.PublicDnsName
			case "PrivateIpAddress":
				resource.ConnectOn = *instance.PrivateIpAddress
			default:
				continue
			}
			resource.Groups = []string{
				*instance.ImageId,
				*instance.InstanceType,
				*instance.VpcId,
			}
			for _, tag := range instance.Tags {
				resource.Groups = append(resource.Groups, strings.ToLower(*tag.Key)+"_"+strings.ToLower(*tag.Value))
			}
			resources = append(resources, resource)
		}
	}
	return resources
}
