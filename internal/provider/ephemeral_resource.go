// Copyright (c) Carlos De La Torre CC-BY-NC-v4 (https://creativecommons.org/licenses/by-nc/4.0/)

package vmworkstation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &VMEphemeralResource{}

func NewVMEphemeralResource() ephemeral.EphemeralResource {
	return &VMEphemeralResource{}
}

type VMEphemeralResource struct {
	// client *http.Client // If applicable, a client can be initialized here.
}

type VMEphemeralResourceModel struct {
	ConfigurableAttribute types.String `tfsdk:"configurable_attribute"`
	Value                 types.String `tfsdk:"value"`
}

func (r *VMEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ephemeral"
}

func (r *VMEphemeralResource) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "In case that we will need to create a Virtual Machine in VmWare Workstation Pro, but is just to make some temporal tasks we can use the Ephemeral Virtual Machine.",
		MarkdownDescription: "Example ephemeral resource",
		Attributes: map[string]schema.Attribute{
			"configurable_attribute": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true, // Ephemeral resources expect their dependencies to already exist.
			},
			"value": schema.StringAttribute{
				Computed: true,
				// Sensitive:           true, // If applicable, mark the attribute as sensitive.
				MarkdownDescription: "Example value",
			},
		},
	}
}

func (r *VMEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data VMEphemeralResourceModel

	// Read Terraform config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
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
	//
	// However, this example hardcodes setting the token attribute to a specific value for brevity.
	data.Value = types.StringValue("token-123")

	// Save data into ephemeral result data
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
