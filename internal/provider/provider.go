// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/option"
)

var _ provider.Provider = &VertexaitxlProvider{}
var _ provider.ProviderWithFunctions = &VertexaitxlProvider{}
var _ provider.ProviderWithEphemeralResources = &VertexaitxlProvider{}

type VertexaitxlProvider struct {
	version string
}

type VertexaitxlProviderModel struct {
	Credentials types.String `tfsdk:"credentials"`
}

func (p *VertexaitxlProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vertexaitxl"
	resp.Version = p.version
}

func (p *VertexaitxlProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"credentials": schema.StringAttribute{
				MarkdownDescription: "Google cloud Platform service account json file path",
				Optional:            true,
			},
		},
	}
}

func (p *VertexaitxlProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data VertexaitxlProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var clientOpts []option.ClientOption

	// Only process credentials if provided
	if !data.Credentials.IsNull() {
		credBytes, err := os.ReadFile(data.Credentials.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read credentials file",
				fmt.Sprintf("Unable to read service account credentials file: %v", err),
			)
			return
		}
		clientOpts = append(clientOpts, option.WithCredentialsJSON(credBytes))
	}

	resp.DataSourceData = clientOpts
	resp.ResourceData = clientOpts
}

func (p *VertexaitxlProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewModelGardenResource,
	}
}

func (p *VertexaitxlProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *VertexaitxlProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *VertexaitxlProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &VertexaitxlProvider{
			version: version,
		}
	}
}
