package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure SSHConfigProvider satisfies various provider interfaces.
var _ provider.Provider = &SSHConfigProvider{}
var _ provider.ProviderWithFunctions = &SSHConfigProvider{}

// SSHConfigProvider defines the provider implementation.
type SSHConfigProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SSHConfigProviderModel describes the provider data model.
type SSHConfigProviderModel struct {
	// No provider-level configuration is needed for this provider
}

func (p *SSHConfigProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sshconfig"
	resp.Version = p.version
}

func (p *SSHConfigProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The SSH Config provider allows you to read SSH configuration files and extract host configurations.",
		MarkdownDescription: "The SSH Config provider allows you to read SSH configuration files and extract host configurations.\n\n" +
			"This provider is useful for reading existing SSH configurations and using them in Terraform configurations.",
	}
}

func (p *SSHConfigProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SSHConfigProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// No configuration is needed for this provider
	resp.DataSourceData = nil
	resp.ResourceData = nil
}

func (p *SSHConfigProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// No resources are provided by this provider
	}
}

func (p *SSHConfigProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewHostDataSource,
	}
}

func (p *SSHConfigProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// No functions are provided by this provider yet
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SSHConfigProvider{
			version: version,
		}
	}
}