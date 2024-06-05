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
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"env": schema.MapAttribute{
										Optional:    true,
										ElementType: types.StringType,
									},
								},
							},
							"custom": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"credentials": schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{
											"username": schema.StringAttribute{
												Required: true,
											},
											"password": schema.StringAttribute{
												Required: true,
											},
										},
									},
									"env": schema.MapAttribute{
										Optional:    true,
										ElementType: types.StringType,
									},
									"health_route": schema.StringAttribute{
										Optional: true,
									},
									"port": schema.Int64Attribute{
										Optional: true,
									},
									"url": schema.StringAttribute{
										Required: true,
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

func clientEndpointToProviderEndpoint(endpoint huggingface.EndpointDetails) endpointResourceModel {
	var image Image
	if endpoint.Model.Image.Huggingface != nil {
		image = Image{
			Huggingface: &Huggingface{
				Env: endpoint.Model.Image.Huggingface.Env,
			},
		}
	} else if endpoint.Model.Image.Custom != nil {
		image = Image{
			Custom: &Custom{
				Env:         endpoint.Model.Image.Custom.Env,
				HealthRoute: endpoint.Model.Image.Custom.HealthRoute,
				Port:        endpoint.Model.Image.Custom.Port,
				URL:         endpoint.Model.Image.Custom.URL,
			},
		}
		if endpoint.Model.Image.Custom.Credentials != nil {
			image.Custom.Credentials = &Credentials{
				Username: endpoint.Model.Image.Custom.Credentials.Username,
				Password: endpoint.Model.Image.Custom.Credentials.Password,
			}
		}
	}

	providerEndpoint := endpointResourceModel{
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
			Framework:  endpoint.Model.Framework,
			Image:      image,
			Repository: endpoint.Model.Repository,
			Revision:   types.StringPointerValue(endpoint.Model.Revision),
			Task:       types.StringPointerValue(endpoint.Model.Task),
		},
		Name: types.StringValue(endpoint.Name),
		Cloud: Cloud{
			Region: endpoint.Provider.Region,
			Vendor: endpoint.Provider.Vendor,
		},
		Type: types.StringValue(endpoint.Type),
	}

	if endpoint.Model.Image.Huggingface != nil && endpoint.Model.Image.Huggingface.Env == nil {
		providerEndpoint.Model.Image.Huggingface.Env = make(map[string]string)
	}

	if endpoint.Model.Image.Custom != nil && endpoint.Model.Image.Custom.Env == nil {
		providerEndpoint.Model.Image.Custom.Env = make(map[string]string)
	}

	return providerEndpoint
}

func providerEndpointToCreateEndpointRequest(endpoint endpointResourceModel) huggingface.CreateEndpointRequest {
	var image huggingface.Image
	if endpoint.Model.Image.Huggingface != nil {
		image = huggingface.Image{
			Huggingface: &huggingface.Huggingface{
				Env: endpoint.Model.Image.Huggingface.Env,
			},
		}
	} else if endpoint.Model.Image.Custom != nil {
		image = huggingface.Image{
			Custom: &huggingface.Custom{
				Env:         endpoint.Model.Image.Custom.Env,
				HealthRoute: endpoint.Model.Image.Custom.HealthRoute,
				Port:        endpoint.Model.Image.Custom.Port,
				URL:         endpoint.Model.Image.Custom.URL,
			},
		}
		if endpoint.Model.Image.Custom.Credentials != nil {
			image.Custom.Credentials = &huggingface.Credentials{
				Username: endpoint.Model.Image.Custom.Credentials.Username,
				Password: endpoint.Model.Image.Custom.Credentials.Password,
			}
		}
	}

	huggingfaceEndpoint := huggingface.CreateEndpointRequest{
		Name:      endpoint.Name.ValueString(),
		AccountId: endpoint.AccountId.ValueStringPointer(),
		Compute: huggingface.Compute{
			Accelerator:  endpoint.Compute.Accelerator,
			InstanceSize: endpoint.Compute.InstanceSize,
			InstanceType: endpoint.Compute.InstanceType,
			Scaling: huggingface.Scaling{
				MaxReplica:         endpoint.Compute.Scaling.MaxReplica,
				MinReplica:         endpoint.Compute.Scaling.MinReplica,
				ScaleToZeroTimeout: endpoint.Compute.Scaling.ScaleToZeroTimeout,
			},
		},
		Model: huggingface.Model{
			Framework:  endpoint.Model.Framework,
			Image:      image,
			Repository: endpoint.Model.Repository,
			Revision:   endpoint.Model.Revision.ValueStringPointer(),
			Task:       endpoint.Model.Task.ValueStringPointer(),
		},
		Provider: huggingface.Provider{
			Region: endpoint.Cloud.Region,
			Vendor: endpoint.Cloud.Vendor,
		},
		Type: endpoint.Type.ValueString(),
	}

	return huggingfaceEndpoint
}

func providerEndpointToUpdateEndpointRequest(endpoint endpointResourceModel) huggingface.UpdateEndpointRequest {
	var image huggingface.Image
	if endpoint.Model.Image.Huggingface != nil {
		image = huggingface.Image{
			Huggingface: &huggingface.Huggingface{
				Env: endpoint.Model.Image.Huggingface.Env,
			},
		}
	} else if endpoint.Model.Image.Custom != nil {
		image = huggingface.Image{
			Custom: &huggingface.Custom{
				Env:         endpoint.Model.Image.Custom.Env,
				HealthRoute: endpoint.Model.Image.Custom.HealthRoute,
				Port:        endpoint.Model.Image.Custom.Port,
				URL:         endpoint.Model.Image.Custom.URL,
			},
		}
		if endpoint.Model.Image.Custom.Credentials != nil {
			image.Custom.Credentials = &huggingface.Credentials{
				Username: endpoint.Model.Image.Custom.Credentials.Username,
				Password: endpoint.Model.Image.Custom.Credentials.Password,
			}
		}
	}

	huggingfaceEndpoint := huggingface.UpdateEndpointRequest{
		Compute: &huggingface.Compute{
			Accelerator:  endpoint.Compute.Accelerator,
			InstanceSize: endpoint.Compute.InstanceSize,
			InstanceType: endpoint.Compute.InstanceType,
			Scaling: huggingface.Scaling{
				MaxReplica:         endpoint.Compute.Scaling.MaxReplica,
				MinReplica:         endpoint.Compute.Scaling.MinReplica,
				ScaleToZeroTimeout: endpoint.Compute.Scaling.ScaleToZeroTimeout,
			},
		},
		Model: &huggingface.Model{
			Framework:  endpoint.Model.Framework,
			Image:      image,
			Repository: endpoint.Model.Repository,
			Revision:   endpoint.Model.Revision.ValueStringPointer(),
			Task:       endpoint.Model.Task.ValueStringPointer(),
		},
		Type: endpoint.Type.ValueStringPointer(),
	}

	return huggingfaceEndpoint
}

func (r *endpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan endpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	existingEndpoints, err := r.client.ListEndpoints()
	if err != nil {
		resp.Diagnostics.AddError(
			"error listing endpoints",
			err.Error(),
		)
		return
	}

	useUpdate := false
	for _, existingEndpoint := range existingEndpoints {
		if existingEndpoint.Name == plan.Name.ValueString() {
			useUpdate = true
			break
		}
	}

	var createdEndpoint huggingface.EndpointDetails

	if useUpdate {
		updateEndpointRequest := providerEndpointToUpdateEndpointRequest(plan)
		createdEndpoint, err = r.client.UpdateEndpoint(plan.Name.ValueString(), updateEndpointRequest)
	} else {
		createEndpointRequest := providerEndpointToCreateEndpointRequest(plan)
		createdEndpoint, err = r.client.CreateEndpoint(createEndpointRequest)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"error creating endpoint",
			err.Error(),
		)
		return
	}

	plan = clientEndpointToProviderEndpoint(createdEndpoint)

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

	state = clientEndpointToProviderEndpoint(endpoint)

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

	endpoint := providerEndpointToUpdateEndpointRequest(plan)

	updatedEndpoint, err := r.client.UpdateEndpoint(plan.Name.ValueString(), endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			"error updating endpoint",
			err.Error(),
		)
		return
	}

	plan = clientEndpointToProviderEndpoint(updatedEndpoint)

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
