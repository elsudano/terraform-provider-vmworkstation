// Copyright (c) Carlos De La Torre CC-BY-NC-v4 (https://creativecommons.org/licenses/by-nc/4.0/)

package vmworkstation

import (
	"context"
	"os"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	HTTPS    types.Bool   `tfsdk:"https"`
	Debug    types.String `tfsdk:"debug"`
}

func (p *VMWProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vmworkstation"
	resp.Version = p.version
}

func (p *VMWProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
A Terraform provider to work with VmWare Workstation Pro, allowing management of virtual machines and other VMware resources.
Supports management through VMware Workstation Pro.
		`,
		MarkdownDescription: `
The VmWare Workstation Pro provider gives Terraform the ability to work with [VmWare Workstation Pro Products][vmware-workstation].
This provider can be used to manage many aspects of a VmWare Workstation Pro environment, including
Virtual Machines, Shared Folders and Networking.

Use the navigation on the left to read about the various resources and data sources supported by this provider.

~> NOTE: This provider requires API REST (vmrest) enable on VmWare Workstation Pro.

[vmware-workstation]: https://www.vmware.com/products/workstation-pro.html
`,
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: `
That's would be the endpoint where we have our VmREST service of VmWare Workstation Pro,
ideally you can configure this value in the provider block of your code, but if you want
You can use the *VMWS_ENDPOINT* environment variable to set this value as well.
`,
				Required:           true,
				Optional:           false,
				Sensitive:          false,
				Description:        "The URL for connect to the API REST.",
				DeprecationMessage: "",
				Validators:         []validator.String{},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: `
That's would be the username that we have configured in our VmREST service of VmWare Workstation Pro,
When you run the vmrest command you can configure the username to access at this API.
You can use the *VMWS_USERNAME* environment variable to set this value.
`, // string,
				Required:           true,                                    // bool,
				Optional:           false,                                   // bool,
				Sensitive:          false,                                   // bool,
				Description:        "The username to use in the API calls.", // string,
				DeprecationMessage: "",                                      // string,
				Validators:         []validator.String{},                    //[]validator.String
			},
			"password": schema.StringAttribute{
				MarkdownDescription: `
That's would be the password that we have configured in our VmREST service of VmWare Workstation Pro,
When you run the vmrest command you can configure the password to access at this API.
You can use the *VMWS_PASSWORD* environment variable to set this value.
`,
				Required:           true,
				Optional:           false,
				Sensitive:          true,
				Description:        "The user password for VMWare Workstation Pro API REST operations.",
				DeprecationMessage: "",
				Validators:         []validator.String{},
			},
			"https": schema.BoolAttribute{
				MarkdownDescription: `
As you konw the VmREST service of VmWare Workstation Pro, can work in HTTP or HTTPS depending if you create
a certificate and run with this certificate the vmrest command or not, for that reason you will need
set at TRUE if you generate the certificate or set FALSE if you want to work in HTTP protocol.
You can use the *VMWS_HTTPS* environment variable to set this value.
`,
				Required:           false,
				Optional:           true,
				Sensitive:          false,
				Description:        "When this have set to true the 'url' connect to over https.",
				DeprecationMessage: "",
			},
			"debug": schema.StringAttribute{
				MarkdownDescription: `
As a piece of software this provider can will be improved, for that reason if you suspect that you found
an issue in the provider, you can enable the debug mode with this field, you have 3 modes NONE, WARN and ERROR.
This levels have been selected from this [link][logging-severity].
You can use the *VMWS_DEBUG* environment variable to set this value.

[logging-severity]: https://developer.hashicorp.com/terraform/plugin/framework/diagnostics#severity
`,
				Required:           false,
				Optional:           true,
				Sensitive:          false,
				Description:        "Enable debug for find errors.",
				DeprecationMessage: "",
				Validators:         []validator.String{},
			},
		},
	}
}

func (p *VMWProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var data VMWProviderModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if data.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown VmWare Workstation API Endpoint",
			"The provider cannot create the VmWare Workstation API client as there is an unknown configuration value for the VmWare Workstation API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VMWS_ENDPOINT environment variable.",
		)
	}
	if data.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown VmWare Workstation API Username",
			"The provider cannot create the VmWare Workstation API client as there is an unknown configuration value for the VmWare Workstation API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VMWS_USERNAME environment variable.",
		)
	}
	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown VmWare Workstation API Password",
			"The provider cannot create the VmWare Workstation API client as there is an unknown configuration value for the VmWare Workstation API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VMWS_PASSWORD environment variable.",
		)
	}
	if data.HTTPS.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("https"),
			"Unknown VmWare Workstation API Insecure",
			"The provider cannot create the VmWare Workstation API client as there is an unknown configuration value for the VmWare Workstation API https param. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VMWS_HTTPS environment variable.",
		)
	}
	if data.Debug.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("debug"),
			"Unknown VmWare Workstation API Debug param.",
			"The provider cannot create the VmWare Workstation API client as there is an unknown configuration value for the VmWare Workstation API Debug param. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VMWS_DEBUG environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	var https bool
	endpoint := os.Getenv("VMWS_ENDPOINT")
	username := os.Getenv("VMWS_USERNAME")
	password := os.Getenv("VMWS_PASSWORD")
	if os.Getenv("VMWS_HTTPS") == "true" {
		https = true
	} else {
		https = false
	}
	debug := os.Getenv("VMWS_DEBUG")
	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}
	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}
	if !data.HTTPS.IsNull() {
		https = data.HTTPS.ValueBool()
	}
	if !data.Debug.IsNull() {
		debug = data.Debug.ValueString()
	}
	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing VmWare Workstation API Endpoint",
			"The provider cannot create the VmWare Workstation API client as there is a missing or empty value for the VmWare Workstation API Endpoint. "+
				"Set the host value in the configuration or use the VMWS_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing VmWare Workstation API Username",
			"The provider cannot create the VmWare Workstation API client as there is a missing or empty value for the VmWare Workstation API username. "+
				"Set the username value in the configuration or use the VMWS_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing VmWare Workstation API Password",
			"The provider cannot create the VmWare Workstation API client as there is a missing or empty value for the VmWare Workstation API password. "+
				"Set the password value in the configuration or use the VMWS_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	// if https {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("https"),
	// 		"Missing VmWare Workstation API Insecure param.",
	// 		"The provider cannot create the VmWare Workstation API client as there is a missing or empty value for the VmWare Workstation API https param. "+
	// 			"Set the https value (true or false) in the configuration or use the VMWS_HTTPS environment variable. "+
	// 			"If either is already set, ensure the value is not empty.",
	// 	)
	// }
	if debug == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("debug"),
			"Missing VmWare Workstation API Debug param.",
			"The provider cannot create the VmWare Workstation API client as there is a missing or empty value for the VmWare Workstation API Debug param. "+
				"Set the Debug value in the configuration or use the VMWS_DEBUG environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	// Create a new VmWare Workstation client using the configuration values
	client := wsapiclient.New()
	err := client.Caller.ConfigClient(endpoint, username, password, https, debug)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create VmWare Workstation API Client",
			"An unexpected error occurred when creating the VmWare Workstation API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"VmWare Workstation Client Error: "+err.Error(),
		)
		return
	}
	// Make the VmWare Workstation client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *VMWProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVMDataSource,
	}
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

func (p *VMWProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
	// return []func() function.Function{
	// 	NewFunction,
	// }
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &VMWProvider{
			version: version,
		}
	}
}
