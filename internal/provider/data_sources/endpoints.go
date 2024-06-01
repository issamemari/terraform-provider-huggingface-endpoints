package data_sources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/issamemari/huggingface-endpoints-client"
)

var (
	_ datasource.DataSource              = &endpointsDataSource{}
	_ datasource.DataSourceWithConfigure = &endpointsDataSource{}
)

func NewEndpointsDataSource() datasource.DataSource {
	return &endpointsDataSource{}
}

type endpointsDataSource struct {
	client *huggingface.Client
}

type endpointsDataSourceModel struct {
	Endpoints []EndpointDetails `tfsdk:"endpoints"`
}

// Configure adds the provider configured client to the data source.
func (d *endpointsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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
	d.client = client
}

func (d *endpointsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoints"
}

func (d *endpointsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoints": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
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
						"status": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"created_at": schema.StringAttribute{
									Computed: true,
								},
								"created_by": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Computed: true,
										},
										"name": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								"error_message": schema.StringAttribute{
									Computed: true,
								},
								"message": schema.StringAttribute{
									Computed: true,
								},
								"private": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"service_name": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								"ready_replica": schema.Int64Attribute{
									Computed: true,
								},
								"state": schema.StringAttribute{
									Computed: true,
								},
								"target_replica": schema.Int64Attribute{
									Computed: true,
								},
								"updated_at": schema.StringAttribute{
									Computed: true,
								},
								"updated_by": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Computed: true,
										},
										"name": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								"url": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *endpointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state endpointsDataSourceModel
	endpoints, err := d.client.ListEndpoints()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Endpoints",
			err.Error(),
		)
		return
	}

	// Map response body to model
	var endpointDetailsList []EndpointDetails

	for _, endpoint := range endpoints {
		endpointDetails := EndpointDetails{
			AccountId: endpoint.AccountId,
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
			Name: endpoint.Name,
			Provider: Provider{
				Region: endpoint.Provider.Region,
				Vendor: endpoint.Provider.Vendor,
			},
			Status: Status{
				CreatedAt:     endpoint.Status.CreatedAt.Format(time.RFC3339),
				CreatedBy:     User{ID: endpoint.Status.CreatedBy.ID, Name: endpoint.Status.CreatedBy.Name},
				ErrorMessage:  endpoint.Status.ErrorMessage,
				Message:       endpoint.Status.Message,
				Private:       Private{ServiceName: endpoint.Status.Private.ServiceName},
				ReadyReplica:  endpoint.Status.ReadyReplica,
				State:         endpoint.Status.State,
				TargetReplica: endpoint.Status.TargetReplica,
				UpdatedAt:     endpoint.Status.UpdatedAt.Format(time.RFC3339),
				UpdatedBy:     User{ID: endpoint.Status.UpdatedBy.ID, Name: endpoint.Status.UpdatedBy.Name},
				URL:           endpoint.Status.URL,
			},
			Type: endpoint.Type,
		}

		endpointDetailsList = append(endpointDetailsList, endpointDetails)
	}

	state.Endpoints = endpointDetailsList

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
