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

type endpointResourceModel struct {
	AccountId types.String `tfsdk:"account_id"`
	Compute   Compute      `tfsdk:"compute"`
	Model     Model        `tfsdk:"model"`
	Name      types.String `tfsdk:"name"`
	Cloud     Cloud        `tfsdk:"cloud"`
	Type      types.String `tfsdk:"type"`
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
								Optional: true,
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
										Optional:    true,
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
						Computed: true,
						Optional: true,
					},
					"task": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"cloud": schema.SingleNestedAttribute{
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
				MaxReplica:         plan.Compute.Scaling.MaxReplica,
				MinReplica:         plan.Compute.Scaling.MinReplica,
				ScaleToZeroTimeout: plan.Compute.Scaling.ScaleToZeroTimeout,
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
			Revision:   plan.Model.Revision.ValueStringPointer(),
			Task:       plan.Model.Task,
		},
		Provider: &huggingface.Provider{
			Region: plan.Cloud.Region,
			Vendor: plan.Cloud.Vendor,
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
		AccountId: types.StringPointerValue(createdEndpoint.AccountId),
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
			Revision:   types.StringPointerValue(createdEndpoint.Model.Revision),
			Task:       createdEndpoint.Model.Task,
		},
		Name: types.StringValue(createdEndpoint.Name),
		Cloud: Cloud{
			Region: createdEndpoint.Provider.Region,
			Vendor: createdEndpoint.Provider.Vendor,
		},
		Type: types.StringValue(createdEndpoint.Type),
	}

	if plan.Model.Image.Huggingface.Env == nil {
		plan.Model.Image.Huggingface.Env = make(map[string]interface{})
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
		AccountId: types.StringPointerValue(endpoint.AccountId),
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
			Revision:   types.StringPointerValue(endpoint.Model.Revision),
			Task:       endpoint.Model.Task,
		},
		Name: types.StringValue(endpoint.Name),
		Cloud: Cloud{
			Region: endpoint.Provider.Region,
			Vendor: endpoint.Provider.Vendor,
		},
		Type: types.StringValue(endpoint.Type),
	}

	if state.Model.Image.Huggingface.Env == nil {
		state.Model.Image.Huggingface.Env = make(map[string]interface{})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *endpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
				MaxReplica:         plan.Compute.Scaling.MaxReplica,
				MinReplica:         plan.Compute.Scaling.MinReplica,
				ScaleToZeroTimeout: plan.Compute.Scaling.ScaleToZeroTimeout,
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
			Revision:   plan.Model.Revision.ValueStringPointer(),
			Task:       plan.Model.Task,
		},
		Provider: &huggingface.Provider{
			Region: plan.Cloud.Region,
			Vendor: plan.Cloud.Vendor,
		},
		Type: plan.Type.ValueString(),
	}

	updatedEndpoint, err := r.client.UpdateEndpoint(endpoint.Name, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			"error updating endpoint",
			err.Error(),
		)
		return
	}

	plan = endpointResourceModel{
		AccountId: types.StringPointerValue(updatedEndpoint.AccountId),
		Compute: Compute{
			Accelerator:  updatedEndpoint.Compute.Accelerator,
			InstanceSize: updatedEndpoint.Compute.InstanceSize,
			InstanceType: updatedEndpoint.Compute.InstanceType,
			Scaling: Scaling{
				MaxReplica:         updatedEndpoint.Compute.Scaling.MaxReplica,
				MinReplica:         updatedEndpoint.Compute.Scaling.MinReplica,
				ScaleToZeroTimeout: updatedEndpoint.Compute.Scaling.ScaleToZeroTimeout,
			},
		},
		Model: Model{
			Framework: updatedEndpoint.Model.Framework,
			Image: Image{
				Huggingface: Huggingface{
					Env: updatedEndpoint.Model.Image.Huggingface.Env,
				},
			},
			Repository: updatedEndpoint.Model.Repository,
			Revision:   types.StringPointerValue(updatedEndpoint.Model.Revision),
			Task:       updatedEndpoint.Model.Task,
		},
		Name: types.StringValue(updatedEndpoint.Name),
		Cloud: Cloud{
			Region: updatedEndpoint.Provider.Region,
			Vendor: updatedEndpoint.Provider.Vendor,
		},
		Type: types.StringValue(updatedEndpoint.Type),
	}

	if updatedEndpoint.Model.Image.Huggingface.Env == nil {
		plan.Model.Image.Huggingface.Env = make(map[string]interface{})
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *endpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state endpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEndpoint(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error deleting endpoint",
			err.Error(),
		)
		return
	}
}
