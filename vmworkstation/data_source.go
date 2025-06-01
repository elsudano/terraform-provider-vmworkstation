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

func NewVMDataSource() datasource.DataSource {
	return &VMDataSource{}
}

type VMDataSource struct {
	client *wsapiclient.Client
}

type VMDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	SourceID types.String `tfsdk:"sourceid"`
	Ip       types.String `tfsdk:"ip"`
}

func (r *VMDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datasource_vm"
}

func (r *VMDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "We can read a VM of VmWare Workstation with this kind of data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				Optional:            false,
				Description:         "That will be the ID of this VM.",
				MarkdownDescription: "When the VM is created the VMWare Workstation Provider assign a new ID at this VM.",
			},
			"sourceid": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the VM that to use for clone at the new.",
				MarkdownDescription: "The VmWare Workstation Provider needs a ID of a created VM to clone this VM in the new one with the new parameters.",
				Validators:          []validator.String{},
			},
			"ip": schema.StringAttribute{
				Computed:            true,
				Description:         "Which is the IP of the instance",
				MarkdownDescription: "When the VM is in PowerON state, we can see which IP have the VM in order to connect with the VM.",
			},
		},
	}
}

func (r *VMDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*wsapiclient.Client)
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
	VM, err := r.client.LoadVM(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read VM, got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("The VM is: %#v", VM))
	data.SourceID = types.StringValue(VM.IdVM)
	data.Ip = types.StringValue("0.0.0.0/0")

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	// data.Id = types.StringValue("545OMDAL1R520604HKNKA6TTK6TBNOHK")
	// data.SourceID = types.StringValue("545OMDAL1R520604HKNKA6TTK6TBNOHK")
	// data.Ip = types.StringValue("0.0.0.0/0")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Info(ctx, "We have read the data source")

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
