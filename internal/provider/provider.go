package provider

import (
	"context"
	"fmt"

	"terraform-provider-huggingface/internal/provider/data_sources"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/issamemari/huggingface-endpoints-client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &huggingfaceProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &huggingfaceProvider{
			version: version,
		}
	}
}

// huggingfaceProvider is the provider implementation.
type huggingfaceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *huggingfaceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "huggingface"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *huggingfaceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"namespace": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

type huggingfaceProviderModel struct {
	Host      types.String `tfsdk:"host"`
	Namespace types.String `tfsdk:"namespace"`
	Token     types.String `tfsdk:"token"`
}

func ValidateConfiguration(config huggingfaceProviderModel, resp *provider.ConfigureResponse) error {
	if config.Host.IsUnknown() || config.Host.IsNull() || config.Host.ValueString() == "" {
		resp.Diagnostics.AddError("host", "HuggingFace API host unknown or empty")
	}
	if config.Namespace.IsUnknown() || config.Namespace.IsNull() || config.Namespace.ValueString() == "" {
		resp.Diagnostics.AddError("namespace", "HuggingFace API namespace unknown or empty")
	}
	if config.Token.IsUnknown() || config.Token.IsNull() || config.Token.ValueString() == "" {
		resp.Diagnostics.AddError("token", "HuggingFace API token unknown or empty")
	}
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("invalid configuration")
	}
	return nil
}

func (p *huggingfaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config huggingfaceProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ValidateConfiguration(config, resp); err != nil {
		return
	}

	host := config.Host.ValueString()
	namespace := config.Namespace.ValueString()
	token := config.Token.ValueString()

	client, err := huggingface.NewClient(&host, &namespace, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create HuggingFace API client",
			err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *huggingfaceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		data_sources.NewEndpointsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *huggingfaceProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
