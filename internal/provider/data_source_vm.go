// Copyright (c) Carlos De La Torre CC-BY-NC-v4 (https://creativecommons.org/licenses/by-nc/4.0/)

package vmworkstation

import (
	"context"
	"fmt"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &VMDataSource{}
var _ datasource.DataSourceWithConfigure = &VMDataSource{}

type VMDataSource struct {
	client *wsapiclient.WSAPIClient
}

type VMDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Denomination types.String `tfsdk:"denomination"`
	Description  types.String `tfsdk:"description"`
	Path         types.String `tfsdk:"path"`
	Processors   types.Int32  `tfsdk:"processors"`
	Memory       types.Int32  `tfsdk:"memory"`
	State        types.String `tfsdk:"state"`
	Ip           types.String `tfsdk:"ip"`
}

func (r *VMDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine"
}

func (r *VMDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
With this item we can read all the properties that we have in our VmWare Workstation Pro folder data,
this means, that we can load a VM in our infrastructure in case that we want to change its properties.
`,
		MarkdownDescription: `
With this item we can read all the properties that we have in our VmWare Workstation Pro folder data,
this means, that we can load a VM of our infrastructure in case that we want to change it properties.

Basically we can read the resources of a VM and then handling this properties with terraform.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "That will be the ID of this VM.",
				MarkdownDescription: `
When the VM is created the VMWare Workstation Provider assign a new ID at this VM.
`,
			},
			"denomination": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VM.",
				MarkdownDescription: `
This will be the name that we can see in the VmWare Workstation.
`,
				Validators: []validator.String{},
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Little bit description of the VM",
				MarkdownDescription: `
Here will have all the description about of the VM, e.g. extra information regarding the purpose of the VM or which is the user and pass of the VM.
`,
			},
			"path": schema.StringAttribute{
				Computed:    true,
				Description: "Absolute path of the VM machine",
				MarkdownDescription: `
Where is the folder where we have the .vmx file of the VM, normally we have this file in the default folder of the VmWare Workstation config.
`,
			},
			"processors": schema.Int32Attribute{
				Computed:    true,
				Description: "Number of processors that will have the VM",
				MarkdownDescription: `
This will be the amount of Processors that the VM will have.
`,
			},
			"memory": schema.Int32Attribute{
				Computed:    true,
				Description: "How much memory will have the VM",
				MarkdownDescription: `
This will be the amount of Memory that the VM will have.
`,
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Which will be the state of the VM when we will deploy it",
				MarkdownDescription: `
That will be state of the VM, that's means that we can have a PowerON VM (on) or a PowerOFF VM (off).
`,
			},
			"ip": schema.StringAttribute{
				Computed:    true,
				Description: "Which is the IP of the instance",
				MarkdownDescription: `
When the VM is in PowerON state, we can see which IP have the VM in order to connect with the VM.
`,
			},
		},
	}
}

func (r *VMDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*wsapiclient.WSAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *wsapiclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *VMDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VMDataSourceModel
	// Read Terraform configuration data into the model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	VM, err := r.client.LoadVMbyName(data.Denomination.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read VM, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("The VM is: %#v", VM))
	if VM.IdVM == "" {
		resp.Diagnostics.AddError(
			"The VM Id field is empty.",
			fmt.Sprintf("Expected string, got: %T. Please report this issue to the provider developers.", VM.IdVM),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("The VM Id field is: %#v", VM.IdVM))
	data.Id = types.StringValue(VM.IdVM)
	if VM.Denomination == "" {
		resp.Diagnostics.AddError(
			"The VM Denomination field is empty.",
			fmt.Sprintf("Expected string, got: %T. Please report this issue to the provider developers.", VM.Denomination),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("The VM Denomination field is: %#v", VM.Denomination))
	data.Denomination = types.StringValue(VM.Denomination)
	if VM.Description == "" {
		resp.Diagnostics.AddError(
			"The VM Description field is empty.",
			fmt.Sprintf("Expected string, got: %T. Please report this issue to the provider developers.", VM.Description),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("The VM Description field is: %#v", VM.Description))
	data.Description = types.StringValue(VM.Description)
	if VM.Path == "" {
		resp.Diagnostics.AddError(
			"The VM Path field is empty.",
			fmt.Sprintf("Expected string, got: %T. Please report this issue to the provider developers.", VM.Path),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("The VM Path field is: %#v", VM.Path))
	data.Path = types.StringValue(VM.Path)
	if VM.CPU.Processors == 0 {
		resp.Diagnostics.AddError(
			"The VM Processors field is empty.",
			fmt.Sprintf("Expected number, got: %T. Please report this issue to the provider developers.", VM.CPU.Processors),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("The VM Processors field is: %#v", VM.CPU.Processors))
	data.Processors = types.Int32Value(VM.CPU.Processors)
	if VM.Memory == 0 {
		resp.Diagnostics.AddError(
			"The VM Memory field is empty.",
			fmt.Sprintf("Expected number, got: %T. Please report this issue to the provider developers.", VM.Memory),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("The VM Memory field is: %#v", VM.Memory))
	data.Memory = types.Int32Value(VM.Memory)
	if VM.PowerStatus == "" {
		resp.Diagnostics.AddError(
			"The VM PowerStatus field is empty.",
			fmt.Sprintf("Expected string, got: %T. Please report this issue to the provider developers.", VM.PowerStatus),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("The VM PowerStatus field is: %#v", VM.PowerStatus))
	data.State = types.StringValue(VM.PowerStatus)
	// data.Ip = types.StringValue(VM.NICS)
	data.Ip = types.StringValue("0.0.0.0/0")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Info(ctx, "We have read the VM")

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func NewVMDataSource() datasource.DataSource {
	return &VMDataSource{}
}
