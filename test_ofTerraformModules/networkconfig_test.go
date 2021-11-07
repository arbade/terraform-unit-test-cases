package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An example of how to test the Terraform module in examples/terraform-aws-network-example using Terratest.
func Test_ShouldBeCreateNetworkConfigs(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.

	var defaultRegion = []string{
		"eu-central-1", //default region
	}

	var restrictedRegionsList = []string{
		"us-east-1",
		"us-east-2",      // Launched 2016
		"us-west-1",      // Launched 2009
		"us-west-2",      // Launched 2011
		"ca-central-1",   // Launched 2016
		"sa-east-1",      // Launched 2011
		"eu-west-1",      // Launched 2007
		"eu-west-2",      // Launched 2016
		"eu-west-3",      // Launched 2017
		"ap-southeast-1", // Launched 2010
		"ap-southeast-2", // Launched 2012
		"ap-northeast-1", // Launched 2011
		"ap-northeast-2", // Launched 2016
		"ap-south-1",     // Launched 2016
		"eu-north-1",     // Launched 2018
	}
	awsRegion := aws.GetRandomStableRegion(t, defaultRegion, restrictedRegionsList)
	//awsRegion := aws.GetRandomRegion(t, defaultRegion, nil)
	//awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// Give the VPC and the subnets correct CIDRs
	vpcCidr := "10.10.0.0/16"
	privateSubnetCidr := "10.10.1.0/24"
	publicSubnetCidr := "10.10.2.0/24"

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/network",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"main_vpc_cidr":       vpcCidr,
			"private_subnet_cidr": privateSubnetCidr,
			"public_subnet_cidr":  publicSubnetCidr,
			"aws_region":          awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	publicSubnetId := terraform.Output(t, terraformOptions, "public_subnet_id")
	privateSubnetId := terraform.Output(t, terraformOptions, "private_subnet_id")
	vpcId := terraform.Output(t, terraformOptions, "main_vpc_id")

	subnets := aws.GetSubnetsForVpc(t, vpcId, awsRegion)

	require.Equal(t, 2, len(subnets))
	// Verify if the network that is supposed to be public is really public
	assert.True(t, aws.IsPublicSubnet(t, publicSubnetId, awsRegion))
	// Verify if the network that is supposed to be private is really private
	assert.False(t, aws.IsPublicSubnet(t, privateSubnetId, awsRegion))
}
