package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccECSTaskExecutionDataSource_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_ecs_task_execution.test"
	clusterName := "aws_ecs_cluster.test"
	taskDefinitionName := "aws_ecs_task_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(t, ecs.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, ecs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskExecutionDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "cluster", clusterName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "task_definition", taskDefinitionName, "arn"),
					resource.TestCheckResourceAttr(dataSourceName, "desired_count", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "launch_type", "FARGATE"),
					resource.TestCheckResourceAttr(dataSourceName, "network_configuration.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "task_arns.#", "1"),
				),
			},
		},
	})
}

func TestAccECSTaskExecutionDataSource_tags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_ecs_task_execution.test"
	clusterName := "aws_ecs_cluster.test"
	taskDefinitionName := "aws_ecs_task_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(t, ecs.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, ecs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskExecutionDataSourceConfig_tags(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "cluster", clusterName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "task_definition", taskDefinitionName, "arn"),
					resource.TestCheckResourceAttr(dataSourceName, "desired_count", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "launch_type", "FARGATE"),
					resource.TestCheckResourceAttr(dataSourceName, "network_configuration.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "task_arns.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.key1", "value1"),
				),
			},
		},
	})
}

func testAccTaskExecutionDataSourceConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    protocol    = "6"
    from_port   = 80
    to_port     = 8000
    cidr_blocks = [aws_vpc.test.cidr_block]
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_ecs_cluster" "test" {
  name = %[1]q
}

resource "aws_ecs_cluster_capacity_providers" "test" {
  cluster_name       = aws_ecs_cluster.test.name
  capacity_providers = ["FARGATE"]
}

resource "aws_ecs_task_definition" "test" {
  family                   = %[1]q
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"

  container_definitions = jsonencode([
    {
      name      = "sleep"
      image     = "busybox"
      cpu       = 10
      command   = ["sleep", "10"]
      memory    = 10
      essential = true
      portMappings = [
        {
          protocol      = "tcp"
          containerPort = 8000
        }
      ]
    }
  ])
}
`, rName)
}

func testAccTaskExecutionDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigVPCWithSubnets(rName, 2),
		testAccTaskExecutionDataSourceConfig_base(rName),
		`
data "aws_ecs_task_execution" "test" {
  depends_on = [aws_ecs_cluster_capacity_providers.test]

  cluster         = aws_ecs_cluster.test.id
  task_definition = aws_ecs_task_definition.test.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.test[*].id
    security_groups  = [aws_security_group.test.id]
    assign_public_ip = false
  }
}
`)
}

func testAccTaskExecutionDataSourceConfig_tags(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(
		acctest.ConfigVPCWithSubnets(rName, 2),
		testAccTaskExecutionDataSourceConfig_base(rName),
		fmt.Sprintf(`
data "aws_ecs_task_execution" "test" {
  depends_on = [aws_ecs_cluster_capacity_providers.test]

  cluster         = aws_ecs_cluster.test.id
  task_definition = aws_ecs_task_definition.test.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.test[*].id
    security_groups  = [aws_security_group.test.id]
    assign_public_ip = false
  }

  tags = {
    %[1]q = %[2]q
  }
}
`, tagKey1, tagValue1))
}
