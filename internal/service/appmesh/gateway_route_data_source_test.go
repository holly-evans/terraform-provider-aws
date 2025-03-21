package appmesh_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/appmesh"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func testAccGatewayRouteDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_appmesh_gateway_route.test"
	dataSourceName := "data.aws_appmesh_gateway_route.test"
	meshName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	vgName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	vsName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	gwRouteName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, appmesh.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, appmesh.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayRouteDataSourceConfig_basic(meshName, vgName, vsName, gwRouteName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "arn", dataSourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "created_date", dataSourceName, "created_date"),
					resource.TestCheckResourceAttrPair(resourceName, "last_updated_date", dataSourceName, "last_updated_date"),
					resource.TestCheckResourceAttrPair(resourceName, "mesh_name", dataSourceName, "mesh_name"),
					resource.TestCheckResourceAttrPair(resourceName, "virtual_router_name", dataSourceName, "virtual_router_name"),
					resource.TestCheckResourceAttrPair(resourceName, "mesh_owner", dataSourceName, "mesh_owner"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "resource_owner", dataSourceName, "resource_owner"),
					resource.TestCheckResourceAttrPair(resourceName, "spec.#", dataSourceName, "spec.#"),
					resource.TestCheckResourceAttrPair(resourceName, "spec.0.grpc_route.#", dataSourceName, "spec.0.grpc_route.#"),
					resource.TestCheckResourceAttrPair(resourceName, "spec.0.http_route.#", dataSourceName, "spec.0.http_route.#"),
					resource.TestCheckResourceAttrPair(resourceName, "spec.0.http2_route.#", dataSourceName, "spec.0.http2_route.#"),
					resource.TestCheckResourceAttrPair(resourceName, "spec.0.priority", dataSourceName, "spec.0.priority"),
					resource.TestCheckResourceAttrPair(resourceName, "tags.%", dataSourceName, "tags.%"),
				),
			},
		},
	})
}

func testAccGatewayRouteDataSourceConfig_basic(meshName, vgName, vsName, gwRouteName string) string {
	return fmt.Sprintf(`
resource "aws_appmesh_mesh" "test" {
  name = %[1]q
}

resource "aws_appmesh_virtual_gateway" "test" {
  name      = %[2]q
  mesh_name = aws_appmesh_mesh.test.id

  spec {
    listener {
      port_mapping {
        port     = 8080
        protocol = "http"
      }
    }
  }
}

resource "aws_appmesh_virtual_service" "test" {
  name      = %[3]q
  mesh_name = aws_appmesh_mesh.test.name

  spec {}
}

resource "aws_appmesh_gateway_route" "test" {
  name                 = %[4]q
  mesh_name            = aws_appmesh_mesh.test.name
  virtual_gateway_name = aws_appmesh_virtual_gateway.test.name

  spec {
    http_route {
      action {
        target {
          virtual_service {
            virtual_service_name = aws_appmesh_virtual_service.test.name
          }
        }
      }

      match {
        prefix = "/"
      }
    }
  }

  tags = {
    Name = %[4]q
  }
}

data "aws_appmesh_gateway_route" "test" {
  name                 = aws_appmesh_gateway_route.test.name
  mesh_name            = aws_appmesh_gateway_route.test.mesh_name
  virtual_gateway_name = aws_appmesh_gateway_route.test.virtual_gateway_name
}
`, meshName, vgName, vsName, gwRouteName)
}
