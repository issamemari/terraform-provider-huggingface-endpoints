package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/issamemari/huggingface-endpoints-client-go"
)

var (
	_ resource.Resource              = &endpointResource{}
	_ resource.ResourceWithConfigure = &endpointResource{}
)

func NewEndpointResource() resource.Resource {
	return &endpointResource{}
}

type endpointResource struct {
	client *huggingface.Client
}

func (r *endpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*huggingface.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"unexpected data source configure type",
			fmt.Sprintf("expected *huggingface.Client, got: %T.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *endpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

func (r *endpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Optional: true,
			},
			"compute": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"accelerator": schema.StringAttribute{
						Required: true,
					},
					"instance_size": schema.StringAttribute{
						Required: true,
					},
					"instance_type": schema.StringAttribute{
						Required: true,
					},
					"scaling": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"max_replica": schema.Int64Attribute{
								Required: true,
							},
							"min_replica": schema.Int64Attribute{
								Required: true,
							},
							"scale_to_zero_timeout": schema.Int64Attribute{
								Required: true,
							},
						},
					},
				},
			},
			"model": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"framework": schema.StringAttribute{
						Required: true,
					},
					"image": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"huggingface": schema.SingleNestedAttribute{
								Required: true,
								Attributes: map[string]schema.Attribute{
									"env": schema.MapAttribute{
										Required:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},
					"repository": schema.StringAttribute{
						Required: true,
					},
					"revision": schema.StringAttribute{
						Required: true,
					},
					"task": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"provider_details": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						Required: true,
					},
					"vendor": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"type": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *endpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan endpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := huggingface.Endpoint{
		Name:      plan.Name.ValueString(),
		AccountId: plan.AccountId.ValueStringPointer(),
		Compute: huggingface.Compute{
			Accelerator:  plan.Compute.Accelerator,
			InstanceSize: plan.Compute.InstanceSize,
			InstanceType: plan.Compute.InstanceType,
			Scaling: huggingface.Scaling{
				MaxReplica:         int(plan.Compute.Scaling.MaxReplica),
				MinReplica:         int(plan.Compute.Scaling.MinReplica),
				ScaleToZeroTimeout: int(plan.Compute.Scaling.ScaleToZeroTimeout),
			},
		},
		Model: huggingface.Model{
			Framework: plan.Model.Framework,
			Image: huggingface.Image{
				Huggingface: huggingface.Huggingface{
					Env: plan.Model.Image.Huggingface.Env,
				},
			},
			Repository: plan.Model.Repository,
			Revision:   plan.Model.Revision,
			Task:       plan.Model.Task,
		},
		Provider: &huggingface.Provider{
			Region: plan.ProviderDetails.Region,
			Vendor: plan.ProviderDetails.Vendor,
		},
		Type: plan.Type.ValueString(),
	}

	createdEndpoint, err := r.client.CreateEndpoint(endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			"error creating endpoint",
			err.Error(),
		)
		return
	}

	plan = endpointResourceModel{
		Compute: Compute{
			Accelerator:  createdEndpoint.Compute.Accelerator,
			InstanceSize: createdEndpoint.Compute.InstanceSize,
			InstanceType: createdEndpoint.Compute.InstanceType,
			Scaling: Scaling{
				MaxReplica:         createdEndpoint.Compute.Scaling.MaxReplica,
				MinReplica:         createdEndpoint.Compute.Scaling.MinReplica,
				ScaleToZeroTimeout: createdEndpoint.Compute.Scaling.ScaleToZeroTimeout,
			},
		},
		Model: Model{
			Framework: createdEndpoint.Model.Framework,
			Image: Image{
				Huggingface: Huggingface{
					Env: createdEndpoint.Model.Image.Huggingface.Env,
				},
			},
			Repository: createdEndpoint.Model.Repository,
			Revision:   createdEndpoint.Model.Revision,
			Task:       createdEndpoint.Model.Task,
		},
		Name: types.StringValue(createdEndpoint.Name),
		ProviderDetails: Provider{
			Region: createdEndpoint.Provider.Region,
			Vendor: createdEndpoint.Provider.Vendor,
		},
		Type: types.StringValue(createdEndpoint.Type),
	}

	if createdEndpoint.AccountId != nil {
		plan.AccountId = types.StringValue(*createdEndpoint.AccountId)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

type endpointResourceModel struct {
	AccountId       types.String `tfsdk:"account_id"`
	Compute         Compute      `tfsdk:"compute"`
	Model           Model        `tfsdk:"model"`
	Name            types.String `tfsdk:"name"`
	ProviderDetails Provider     `tfsdk:"provider_details"`
	Type            types.String `tfsdk:"type"`
}

func (r *endpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state endpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint, err := r.client.GetEndpoint(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading endpoint",
			"could not read endpoint named "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	state = endpointResourceModel{
		Compute: Compute{
			Accelerator:  endpoint.Compute.Accelerator,
			InstanceSize: endpoint.Compute.InstanceSize,
			InstanceType: endpoint.Compute.InstanceType,
			Scaling: Scaling{
				MaxReplica:         endpoint.Compute.Scaling.MaxReplica,
				MinReplica:         endpoint.Compute.Scaling.MinReplica,
				ScaleToZeroTimeout: endpoint.Compute.Scaling.ScaleToZeroTimeout,
			},
		},
		Model: Model{
			Framework: endpoint.Model.Framework,
			Image: Image{
				Huggingface: Huggingface{
					Env: endpoint.Model.Image.Huggingface.Env,
				},
			},
			Repository: endpoint.Model.Repository,
			Revision:   endpoint.Model.Revision,
			Task:       endpoint.Model.Task,
		},
		Name: types.StringValue(endpoint.Name),
		ProviderDetails: Provider{
			Region: endpoint.Provider.Region,
			Vendor: endpoint.Provider.Vendor,
		},
		Type: types.StringValue(endpoint.Type),
	}

	if endpoint.AccountId != nil {
		state.AccountId = types.StringValue(*endpoint.AccountId)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *endpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *endpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
