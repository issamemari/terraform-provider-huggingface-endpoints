package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/issamemari/huggingface-endpoints-client"
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
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *huggingface.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
				Computed: true,
			},
			"compute": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"accelerator": schema.StringAttribute{
						Computed: true,
					},
					"instance_size": schema.StringAttribute{
						Computed: true,
					},
					"instance_type": schema.StringAttribute{
						Computed: true,
					},
					"scaling": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"max_replica": schema.Int64Attribute{
								Computed: true,
							},
							"min_replica": schema.Int64Attribute{
								Computed: true,
							},
							"scale_to_zero_timeout": schema.Int64Attribute{
								Computed: true,
							},
						},
					},
				},
			},
			"model": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"framework": schema.StringAttribute{
						Computed: true,
					},
					"image": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"huggingface": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"env": schema.MapAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},
					"repository": schema.StringAttribute{
						Computed: true,
					},
					"revision": schema.StringAttribute{
						Computed: true,
					},
					"task": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"provider": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						Computed: true,
					},
					"vendor": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

type endpointResourceModel struct {
	Endpoint Endpoint `tfsdk:"endpoint"`
}

func (r *endpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan endpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var items []huggingface.Endpoint
	endpoint =huggingface.Endpoint{
			AccountId: plan.Endpoint.AccountId,
			Compute: huggingface.Compute{
				Accelerator:  plan.Endpoint.Compute.Accelerator,
				InstanceSize: plan.Endpoint.Compute.InstanceSize,
				InstanceType: plan.Endpoint.Compute.InstanceType,
				Scaling: huggingface.Scaling{
					MaxReplica:         int(plan.Endpoint.Compute.Scaling.MaxReplica),
					MinReplica:         int(plan.Endpoint.Compute.Scaling.MinReplica),
					ScaleToZeroTimeout: int(plan.Endpoint.Compute.Scaling.ScaleToZeroTimeout),
				},
			},
			Model: huggingface.Model{
				Framework: plan.Endpoint.Model.Framework,
				Image: huggingface.Image{
					Huggingface: huggingface.Huggingface{
						Env: plan.Endpoint.Model.Image.Huggingface.Env,
					},
				},
				Repository: plan.Endpoint.Model.Repository,
				Revision:   plan.Endpoint.Model.Revision,
				Task:       plan.Endpoint.Model.Task,
			},
			Name: plan.Endpoint.Name,
			Provider: huggingface.Provider{
				Region: plan.Endpoint.Provider.Region,
				Vendor: plan.Endpoint.Provider.Vendor,
			},
		}
	}

	endpointCreationResponse, err := r.client.CreateEndpoint(endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *endpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *endpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *endpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
