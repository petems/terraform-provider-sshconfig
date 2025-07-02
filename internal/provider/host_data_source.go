package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/petems/terraform-provider-sshconfig/internal/sshconfig"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &HostDataSource{}

func NewHostDataSource() datasource.DataSource {
	return &HostDataSource{}
}

// HostDataSource defines the data source implementation.
type HostDataSource struct{}

// HostDataSourceModel describes the data source data model.
type HostDataSourceModel struct {
	Path     types.String `tfsdk:"path"`
	Host     types.String `tfsdk:"host"`
	Rendered types.String `tfsdk:"rendered"`
	HostMap  types.Map    `tfsdk:"host_map"`
	ID       types.String `tfsdk:"id"`
}

func (d *HostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (d *HostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSH host configuration data source.\n\n" +
			"This data source allows you to read SSH configuration files and extract host-specific configurations.",
		Description: "SSH host configuration data source",

		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				MarkdownDescription: "The path to the SSH config file you want to read. Defaults to `/etc/ssh/ssh_config`.",
				Description:         "The path to the SSH config file you want to read",
				Optional:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host you want to lookup in the SSH config file.",
				Description:         "The host you want to lookup",
				Required:            true,
			},
			"rendered": schema.StringAttribute{
				MarkdownDescription: "The host information rendered as a multi-line string.",
				Description:         "The host information rendered as a multi-line string",
				Computed:            true,
			},
			"host_map": schema.MapAttribute{
				MarkdownDescription: "The host information as a map of configuration keys to values.",
				Description:         "The host information as a map",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this data source",
				Description:         "Unique identifier for this data source",
				Computed:            true,
			},
		},
	}
}

func (d *HostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HostDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set default path if not provided
	configPath := "/etc/ssh/ssh_config"
	if !data.Path.IsNull() {
		configPath = data.Path.ValueString()
	}

	host := data.Host.ValueString()

	// Parse the SSH config file using our internal parser
	config, err := sshconfig.ParseFile(configPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse SSH Config File",
			fmt.Sprintf("Unable to parse SSH config file at %s: %s", configPath, err),
		)
		return
	}

	// Find the host configuration
	hostConfig := config.FindHost(host)
	if hostConfig == nil {
		resp.Diagnostics.AddError(
			"Host Not Found",
			fmt.Sprintf("Could not find host %s in SSH config file %s", host, configPath),
		)
		return
	}

	// Get the rendered string representation
	hostDeclaration := hostConfig.String()

	// Get merged options (this applies SSH precedence rules)
	mergedOptions := config.GetMergedOptions(host)

	// Convert the map to types.Map
	hostMapValue, diags := types.MapValueFrom(ctx, types.StringType, mergedOptions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the data source values
	data.Path = types.StringValue(configPath)
	data.Host = types.StringValue(host)
	data.Rendered = types.StringValue(hostDeclaration)
	data.HostMap = hostMapValue
	data.ID = types.StringValue(fmt.Sprintf("%s:%s", configPath, host))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
