package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type Endpoint struct {
	AccountId *string `tfsdk:"account_id"`
	Compute   Compute `tfsdk:"compute"`
	Model     Model   `tfsdk:"model"`
	Name      string  `tfsdk:"name"`
	Cloud     Cloud   `tfsdk:"cloud"`
	Status    *Status `tfsdk:"status"`
	Type      string  `tfsdk:"type"`
}

type Compute struct {
	Accelerator  string  `tfsdk:"accelerator"`
	InstanceSize string  `tfsdk:"instance_size"`
	InstanceType string  `tfsdk:"instance_type"`
	Scaling      Scaling `tfsdk:"scaling"`
}

type Scaling struct {
	MaxReplica         int  `tfsdk:"max_replica"`
	MinReplica         int  `tfsdk:"min_replica"`
	ScaleToZeroTimeout *int `tfsdk:"scale_to_zero_timeout"`
}

type Model struct {
	Framework  string       `tfsdk:"framework"`
	Image      Image        `tfsdk:"image"`
	Repository string       `tfsdk:"repository"`
	Revision   types.String `tfsdk:"revision"`
	Task       types.String `tfsdk:"task"`
}

type Image struct {
	Huggingface *Huggingface `tfsdk:"huggingface"`
	Custom      *Custom      `tfsdk:"custom"`
}

type Custom struct {
	Credentials *Credentials      `tfsdk:"credentials"`
	Env         map[string]string `tfsdk:"env"`
	HealthRoute *string           `tfsdk:"health_route"`
	Port        *int              `tfsdk:"port"`
	URL         string            `tfsdk:"url"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Huggingface struct {
	Env map[string]string `tfsdk:"env"`
}

type Cloud struct {
	Region string `tfsdk:"region"`
	Vendor string `tfsdk:"vendor"`
}

type Status struct {
	CreatedAt     string  `tfsdk:"created_at"`
	CreatedBy     User    `tfsdk:"created_by"`
	ErrorMessage  string  `tfsdk:"error_message"`
	Message       string  `tfsdk:"message"`
	Private       Private `tfsdk:"private"`
	ReadyReplica  int     `tfsdk:"ready_replica"`
	State         string  `tfsdk:"state"`
	TargetReplica int     `tfsdk:"target_replica"`
	UpdatedAt     string  `tfsdk:"updated_at"`
	UpdatedBy     User    `tfsdk:"updated_by"`
	URL           string  `tfsdk:"url"`
}

type User struct {
	ID   string `tfsdk:"id"`
	Name string `tfsdk:"name"`
}

type Private struct {
	ServiceName string `tfsdk:"service_name"`
}
