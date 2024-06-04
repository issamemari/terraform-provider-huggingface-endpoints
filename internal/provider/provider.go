package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/issamemari/huggingface-endpoints-client-go"
)

var (
	_ provider.Provider = &huggingfaceProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &huggingfaceProvider{
			version: version,
		}
	}
}

type huggingfaceProvider struct {
	version string
}

// Metadata returns the provider type name.
func (p *huggingfaceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "huggingface-endpoints"
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
		resp.Diagnostics.AddError("host", "huggingface api host unknown or empty")
	}
	if config.Namespace.IsUnknown() || config.Namespace.IsNull() || config.Namespace.ValueString() == "" {
		resp.Diagnostics.AddError("namespace", "huggingface api namespace unknown or empty")
	}
	if config.Token.IsUnknown() || config.Token.IsNull() || config.Token.ValueString() == "" {
		resp.Diagnostics.AddError("token", "huggingface api token unknown or empty")
	}
	if resp.Diagnostics.HasError() {
		return fmt.Errorf("invalid configuration")
	}
	return nil
}

func (p *huggingfaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "configuring huggingface provider")

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

	ctx = tflog.SetField(ctx, "huggingface_host", host)
	ctx = tflog.SetField(ctx, "huggingface_namespace", namespace)
	ctx = tflog.SetField(ctx, "huggingface_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "huggingface_token")
	tflog.Debug(ctx, "creating huggingface client")

	client, err := huggingface.NewClient(&host, &namespace, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to create huggingface api client",
			err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "huggingface provider configured", map[string]any{"success": true})
}

func (p *huggingfaceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *huggingfaceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEndpointResource,
	}
}
