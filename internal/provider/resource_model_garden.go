// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"google.golang.org/api/option"
)

func stringToGenAIType(strType string) genai.Type {
	switch strings.ToLower(strType) {
	case "string":
		return genai.TypeString
	case "number":
		return genai.TypeNumber
	case "integer":
		return genai.TypeInteger
	case "bool":
		return genai.TypeBoolean
	case "array":
		return genai.TypeArray
	case "object":
		return genai.TypeObject
	default:
		return genai.TypeString
	}
}

type Schema struct {
	// Optional. The type of the data.
	Type string
	// Optional. The format of the data.
	// Supported formats:
	//
	//	for NUMBER type: "float", "double"
	//	for INTEGER type: "int32", "int64"
	//	for STRING type: "email", "byte", etc
	Format string
	// Optional. The title of the Schema.
	Title string
	// Optional. The description of the data.
	Description string
	// Optional. Indicates if the value may be null.
	Nullable bool
	// Optional. SCHEMA FIELDS FOR TYPE ARRAY
	// Schema of the elements of Type.ARRAY.
	Items *Schema
	// Optional. Minimum number of the elements for Type.ARRAY.
	MinItems int64
	// Optional. Maximum number of the elements for Type.ARRAY.
	MaxItems int64
	// Optional. Possible values of the element of Type.STRING with enum format.
	// For example we can define an Enum Direction as :
	// {type:STRING, format:enum, enum:["EAST", NORTH", "SOUTH", "WEST"]}
	Enum []string
	// Optional. SCHEMA FIELDS FOR TYPE OBJECT
	// Properties of Type.OBJECT.
	Properties map[string]*Schema
	// Optional. Required properties of Type.OBJECT.
	Required []string
	// Optional. Minimum number of the properties for Type.OBJECT.
	MinProperties int64
	// Optional. Maximum number of the properties for Type.OBJECT.
	MaxProperties int64
	// Optional. SCHEMA FIELDS FOR TYPE INTEGER and NUMBER
	// Minimum value of the Type.INTEGER and Type.NUMBER
	Minimum float64
	// Optional. Maximum value of the Type.INTEGER and Type.NUMBER
	Maximum float64
	// Optional. SCHEMA FIELDS FOR TYPE STRING
	// Minimum length of the Type.STRING
	MinLength int64
	// Optional. Maximum length of the Type.STRING
	MaxLength int64
	// Optional. Pattern of the Type.STRING to restrict a string to a regular
	// expression.
	Pattern string
}

func convertToGenAISchema(config *Schema) *genai.Schema {
	if config == nil {
		return nil
	}

	schema := &genai.Schema{
		Type:          stringToGenAIType(config.Type),
		Format:        config.Format,
		Title:         config.Title,
		Description:   config.Description,
		Nullable:      config.Nullable,
		MinItems:      config.MinItems,
		MaxItems:      config.MaxItems,
		Enum:          config.Enum,
		Required:      config.Required,
		MinProperties: config.MinProperties,
		MaxProperties: config.MaxProperties,
		Minimum:       config.Minimum,
		Maximum:       config.Maximum,
		MinLength:     config.MinLength,
		MaxLength:     config.MaxLength,
		Pattern:       config.Pattern,
	}
	if len(config.Properties) > 0 {
		schema.Properties = make(map[string]*genai.Schema)
		for propertykey, propertyvalue := range config.Properties {
			schema.Properties[propertykey] = convertToGenAISchema(propertyvalue)
		}
	}
	if config.Items != nil {
		schema.Items = convertToGenAISchema(config.Items)
	}
	return schema
}

var _ resource.Resource = &ModelGardenResource{}

func NewModelGardenResource() resource.Resource {
	return &ModelGardenResource{}
}

type ModelGardenResource struct {
	clientOptions []option.ClientOption
}

type ModelGardenResourceModel struct {
	Prompt         types.String `tfsdk:"prompt"`
	ProjectID      types.String `tfsdk:"project_id"`
	Location       types.String `tfsdk:"location"`
	ModelName      types.String `tfsdk:"model_name"`
	ResponseSchema types.String `tfsdk:"response_schema"`
	Response       types.String `tfsdk:"response"`
}

func (r *ModelGardenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model_garden"
}

func (r *ModelGardenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Model Garden Resource",

		Attributes: map[string]schema.Attribute{
			"prompt": schema.StringAttribute{
				MarkdownDescription: "The prompt to send to the model",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"response_schema": schema.StringAttribute{
				MarkdownDescription: "The Response Schema for output response",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Google Cloud Project ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Location for the VertexAI API",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"model_name": schema.StringAttribute{
				MarkdownDescription: "Name of the model to use",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"response": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Response from the model",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ModelGardenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientOpts, err := req.ProviderData.([]option.ClientOption)
	if !err {
		resp.Diagnostics.AddError(
			"Error creating Vertex AI client",
			fmt.Sprintf("Could not create client: %v", err),
		)
		return
	}
	r.clientOptions = clientOpts
}

func (r *ModelGardenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ModelGardenResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := genai.NewClient(ctx, plan.ProjectID.ValueString(), plan.Location.ValueString(), r.clientOptions...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Vertex AI client",
			fmt.Sprintf("Could not create client: %v", err),
		)
		return
	}
	defer client.Close()

	model := client.GenerativeModel(plan.ModelName.ValueString())

	model.GenerationConfig.ResponseMIMEType = "application/json"

	var schemaMap Schema

	if err := json.Unmarshal([]byte(plan.ResponseSchema.ValueString()), &schemaMap); err != nil {
		resp.Diagnostics.AddError(
			"Error while reading the json string of response schema",
			fmt.Sprintf("Could not create client: %v", err),
		)
		return
	}

	generatedSchema := convertToGenAISchema(&schemaMap)

	model.GenerationConfig.ResponseSchema = generatedSchema

	result, err := model.GenerateContent(ctx, genai.Text(plan.Prompt.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating content",
			fmt.Sprintf("Could not generate content: %v", err),
		)
		return
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		resp.Diagnostics.AddError(
			"Empty response",
			"The model returned an empty response",
		)
		return
	}

	plan.Response = types.StringValue(fmt.Sprintf("%v", result.Candidates[0].Content.Parts[0]))

	tflog.Trace(ctx, "created a resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ModelGardenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ModelGardenResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ModelGardenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ModelGardenResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := genai.NewClient(ctx, plan.ProjectID.ValueString(), plan.Location.ValueString(), r.clientOptions...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Vertex AI client",
			fmt.Sprintf("Could not create client: %v", err),
		)
		return
	}
	defer client.Close()

	model := client.GenerativeModel(plan.ModelName.ValueString())

	model.GenerationConfig.ResponseMIMEType = "application/json"

	var schemaMap Schema

	if err := json.Unmarshal([]byte(plan.ResponseSchema.ValueString()), &schemaMap); err != nil {
		resp.Diagnostics.AddError(
			"Error while reading the json string of response schema",
			fmt.Sprintf("Could not create client: %v", err),
		)
		return
	}

	generatedSchema := convertToGenAISchema(&schemaMap)

	model.GenerationConfig.ResponseSchema = generatedSchema

	result, err := model.GenerateContent(ctx, genai.Text(plan.Prompt.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating content",
			fmt.Sprintf("Could not generate content: %v", err),
		)
		return
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		resp.Diagnostics.AddError(
			"Empty response",
			"The model returned an empty response",
		)
		return
	}

	plan.Response = types.StringValue(fmt.Sprintf("%v", result.Candidates[0].Content.Parts[0]))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ModelGardenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}
