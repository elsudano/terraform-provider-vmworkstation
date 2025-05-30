package vmworkstation

import (
	"context"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &VMWProvider{}
var _ provider.ProviderWithFunctions = &VMWProvider{}
var _ provider.ProviderWithEphemeralResources = &VMWProvider{}

type VMWProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type VMWProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	HTTPS    types.Bool   `tfsdk:"https"`
	Debug    types.String `tfsdk:"debug"`
}

func (p *VMWProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vmws"
	resp.Version = p.version
}

func (p *VMWProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The URL for connect to the API REST.",
				Required:            true,
				Optional:            false,
				Sensitive:           false,
				Description:         "The URL for connect to the API REST.",
				DeprecationMessage:  "",
				Validators:          []validator.String{},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The user to use in the API calls.", // string,
				Required:            true,                                // bool,
				Optional:            false,                               // bool,
				Sensitive:           false,                               // bool,
				Description:         "The user to use in the API calls.", // string,
				DeprecationMessage:  "",                                  // string,
				Validators:          []validator.String{},                //[]validator.String
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The user password for VMWare Workstation Pro API REST operations.",
				Required:            true,
				Optional:            false,
				Sensitive:           true,
				Description:         "The user password for VMWare Workstation Pro API REST operations.",
				DeprecationMessage:  "",
				Validators:          []validator.String{},
			},
			"https": schema.BoolAttribute{
				MarkdownDescription: "When this have set to true the 'url' connect to over https.",
				Required:            false,
				Optional:            true,
				Sensitive:           false,
				Description:         "When this have set to true the 'url' connect to over https.",
				DeprecationMessage:  "",
			},
			"debug": schema.StringAttribute{
				MarkdownDescription: "Enable debug for find errors.",
				Required:            false,
				Optional:            true,
				Sensitive:           false,
				Description:         "Enable debug for find errors.",
				DeprecationMessage:  "",
				Validators:          []validator.String{},
			},
		},
	}
}

func (p *VMWProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data VMWProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := wsapiclient.NewClient(data.Endpoint.String(), data.User.String(), data.Password.String(), data.HTTPS.ValueBool(), data.Debug.String())
	if err != nil {
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *VMWProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVMResource,
	}
}

func (p *VMWProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewVMEphemeralResource,
	}
}

func (p *VMWProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVMDataSource,
	}
}

func (p *VMWProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewVMFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &VMWProvider{
			version: version,
		}
	}
}
