package awsinventory

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pythonandchips/azad/plugin"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
	"github.com/stretchr/testify/assert"
)

func TestListEC2Resources(t *testing.T) {
	instance := &ec2.Instance{
		PublicIpAddress:  aws.String("192.186.0.1"),
		PrivateDnsName:   aws.String("aws.internal.ip"),
		PublicDnsName:    aws.String("aws.public.ip"),
		PrivateIpAddress: aws.String("10.0.0.1"),
		ImageId:          aws.String("ami-1234567"),
		InstanceType:     aws.String("m3.medium"),
		VpcId:            aws.String("vpc-13453"),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Environment"),
				Value: aws.String("Development"),
			},
		},
	}
	fakeEC2 := &fakeEC2{
		instances: []*ec2.Instance{instance},
	}
	ec2Session = func(context plugin.Context) ec2iface {
		return fakeEC2
	}
	t.Run("with basic configuration", func(t *testing.T) {
		context := plugintesting.NewFakeContext()

		resources, err := ec2Resources(context)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		t.Run("filter for only running instances", func(t *testing.T) {
			filters := fakeEC2.input.Filters

			if len(filters) != 1 {
				t.Fatalf("expected %d filters but got %d", 1, len(filters))
			}
			filter := filters[0]
			assert.Equal(t, filter.Name, aws.String("instance-state-name"))
			assert.Equal(t, filter.Values, aws.StringSlice([]string{"running"}))
		})
		t.Run("returns the instances as resources", func(t *testing.T) {
			if len(resources) != 1 {
				t.Fatalf("expected %d instances but got %d", 1, len(resources))
			}
			resource := resources[0]
			assert.Equal(t, resource.ConnectOn, "192.186.0.1")

			t.Run("configures the instance with groups", func(t *testing.T) {
				groups := resource.Groups
				testContains(t, groups, "ami-1234567")
				testContains(t, groups, "m3.medium")
				testContains(t, groups, "vpc-13453")
				testContains(t, groups, "environment_development")
			})
		})
	})
	t.Run("with connect_on as PrivateDnsName", func(t *testing.T) {
		context := plugin.NewInventoryContext(map[string]string{
			"connect_on": "PrivateDnsName",
		})

		resources, err := ec2Resources(context)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		assert.Equal(t, resources[0].ConnectOn, "aws.internal.ip")
	})
	t.Run("with connect_on as PublicDnsName", func(t *testing.T) {
		context := plugin.NewInventoryContext(map[string]string{
			"connect_on": "PublicDnsName",
		})

		resources, err := ec2Resources(context)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		assert.Equal(t, resources[0].ConnectOn, "aws.public.ip")
	})
	t.Run("with connect_on as PrivateIpAddress", func(t *testing.T) {
		context := plugin.NewInventoryContext(map[string]string{
			"connect_on": "PrivateIpAddress",
		})

		resources, err := ec2Resources(context)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		assert.Equal(t, resources[0].ConnectOn, "10.0.0.1")
	})
	t.Run("with unrecognized connect_on ignores servers", func(t *testing.T) {
		context := plugin.NewInventoryContext(map[string]string{
			"connect_on": "NotRealValue",
		})

		resources, err := ec2Resources(context)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(resources) != 0 {
			t.Errorf("expected %d resource but got %d", 0, len(resources))
		}
	})
}

func testContains(t *testing.T, groups []string, searchString string) {
	for _, str := range groups {
		if str == searchString {
			return
		}
	}
	t.Errorf("expected slice to contain %s", searchString)
}

type fakeEC2 struct {
	input     *ec2.DescribeInstancesInput
	instances []*ec2.Instance
}

func (fakeEC2 *fakeEC2) DescribeInstancesPages(input *ec2.DescribeInstancesInput,
	pager func(*ec2.DescribeInstancesOutput, bool) bool) error {
	fakeEC2.input = input

	result := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{Instances: fakeEC2.instances},
		},
	}

	pager(result, true)

	return nil
}
