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

func NewVMResource() resource.Resource {
	return &VMResource{}
}

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
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }
	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue("MYNEWVMISCREATED")
	data.SourceID = types.StringValue("545OMDAL1R520604HKNKA6TTK6TBNOHK")
	data.Ip = types.StringValue("0.0.0.0/0")
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VMResourceModel
	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }
	// Save updated data into Terraform state
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
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *VMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
