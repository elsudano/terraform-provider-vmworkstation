// Copyright (c) Carlos De La Torre CC-BY-NC-v4 (https://creativecommons.org/licenses/by-nc/4.0/)

package vmworkstation

import (
	"context"
	"fmt"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &VMResource{}
var _ resource.ResourceWithImportState = &VMResource{}

type VMResource struct {
	client *wsapiclient.WSAPIClient
}

// ExampleResourceModel describes the resource data model.
type VMResourceModel struct {
	Id           types.String `tfsdk:"id"`
	SourceID     types.String `tfsdk:"sourceid"`
	Denomination types.String `tfsdk:"denomination"`
	Description  types.String `tfsdk:"description"`
	Path         types.String `tfsdk:"path"`
	Processors   types.Int32  `tfsdk:"processors"`
	Memory       types.Int32  `tfsdk:"memory"`
	State        types.String `tfsdk:"state"`
	Ip           types.String `tfsdk:"ip"`
}

func (r *VMResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine"
}

func (r *VMResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
		Description: `
The principal resource of VmWare Workstation Pro is a Virtual Machine, with this resource we can create a
clone from a VM that we had in our VmWare Workstation Pro folder.
`,
		MarkdownDescription: `
We can create a VM within of VmWare Workstation with this kind of resource.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "That will be the ID of this VM.",
				MarkdownDescription: `
When the VM is created the VMWare Workstation Provider assign a new ID at this VM.
`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sourceid": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "The ID of the VM that to use for clone at the new.",
				MarkdownDescription: `
The VmWare Workstation Provider needs a ID of a created VM to clone this VM in the new one with the new parameters.
`,
				Validators: []validator.String{},
			},
			"denomination": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "The name of the VM.",
				MarkdownDescription: `
This will be the name that we can see in the VmWare Workstation.
`,
			},
			"description": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "Little bit description of the VM",
				MarkdownDescription: `
Here will have all the description about of the VM, e.g. extra information regarding the purpose of the VM or which is the user and pass of the VM.
`,
			},
			"path": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Absolute path of the VM machine",
				MarkdownDescription: `
Where is the folder where we have the .vmx file of the VM, normally we have this file in the default folder of the VmWare Workstation config.
`,
			},
			// "image": schema.StringAttribute{
			// 	Required:            true,
			// 	Optional:            false,
			// 	Computed:            false,
			// 	Sensitive:           false,
			// 	Description:         "String with ID for the image that create the VM.",
			// 	MarkdownDescription: `String with ID for the image that create the VM.`,
			// 	DeprecationMessage:  "",
			// 	Validators:          []validator.String{},
			// 	PlanModifiers: []planmodifier.String{
			// 		stringplanmodifier.UseStateForUnknown(),
			// 	},
			// 	Default: stringdefault.StaticString("example value when not configured"),
			// },
			"processors": schema.Int32Attribute{
				Required:    true,
				Optional:    false,
				Description: "Number of processors that will have the VM",
				MarkdownDescription: `
This will be the amount of Processors that the VM will have.
`,
			},
			"memory": schema.Int32Attribute{
				Required:    true,
				Optional:    false,
				Description: "How much memory will have the VM",
				MarkdownDescription: `
This will be the amount of Memory that the VM will have.
`,
			},
			"state": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "Which will be the state of the VM when we will deploy it",
				MarkdownDescription: `
That will be state of the VM, that's means that we can have a PowerON VM (on) or a PowerOFF VM (off).
`,
			},
			"ip": schema.StringAttribute{
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: "Which is the IP of the instance",
				MarkdownDescription: `
When the VM is in PowerON state, we can see which IP have the VM in order to connect with the VM.
`,
				Default: stringdefault.StaticString("0.0.0.0/0"),
			},
		},
	}
}

func (r *VMResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*wsapiclient.WSAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *wsapiclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *VMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VMResourceModel
	// Read Terraform plan data into the model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	VM, err := r.client.CreateVM(data.SourceID.ValueString(), data.Denomination.ValueString(), data.Description.ValueString(), int(data.Processors.ValueInt32()), int(data.Memory.ValueInt32()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create VM, got error: %s", err))
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
	// You will need fix that regarding of the network
	data.Ip = types.StringValue("0.0.0.0/0")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Info(ctx, "We have Created a new VM")

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *VMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VMDataSourceModel
	// Read Terraform configuration data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	VM, err := r.client.LoadVM(data.Id.ValueString())
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
	// You will need fix that regarding of the network
	data.Ip = types.StringValue("0.0.0.0/0")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Info(ctx, "We have Read the VM")

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *VMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VMResourceModel
	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	VM, err := r.client.LoadVM(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read VM, got error: %s", err))
		return
	}
	err = r.client.UpdateVM(VM, data.Denomination.ValueString(), data.Description.ValueString(), int(data.Processors.ValueInt32()), int(data.Memory.ValueInt32()), data.State.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update VM, got error: %s", err))
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
	// You will need fix that regarding of the network
	data.Ip = types.StringValue("0.0.0.0/0")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Info(ctx, "We have Updated the VM")

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *VMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VMResourceModel
	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	VM, err := r.client.LoadVM(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Read VM, got error: %s", err))
		return
	}
	err = r.client.DeleteVM(VM)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete VM, got error: %s", err))
		return
	}
	tflog.Info(ctx, "We have Deleted the VM")
}

func (r *VMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func NewVMResource() resource.Resource {
	return &VMResource{}
}
