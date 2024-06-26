// Code generated by smithy-go-codegen DO NOT EDIT.

package neptune

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	internalauth "github.com/aws/aws-sdk-go-v2/internal/auth"
	presignedurlcust "github.com/aws/aws-sdk-go-v2/service/internal/presigned-url"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates a new Amazon Neptune DB cluster. You can use the
// ReplicationSourceIdentifier parameter to create the DB cluster as a Read Replica
// of another DB cluster or Amazon Neptune DB instance. Note that when you create a
// new cluster using CreateDBCluster directly, deletion protection is disabled by
// default (when you create a new production cluster in the console, deletion
// protection is enabled by default). You can only delete a DB cluster if its
// DeletionProtection field is set to false .
func (c *Client) CreateDBCluster(ctx context.Context, params *CreateDBClusterInput, optFns ...func(*Options)) (*CreateDBClusterOutput, error) {
	if params == nil {
		params = &CreateDBClusterInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateDBCluster", params, optFns, c.addOperationCreateDBClusterMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateDBClusterOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateDBClusterInput struct {

	// The DB cluster identifier. This parameter is stored as a lowercase string.
	// Constraints:
	//   - Must contain from 1 to 63 letters, numbers, or hyphens.
	//   - First character must be a letter.
	//   - Cannot end with a hyphen or contain two consecutive hyphens.
	// Example: my-cluster1
	//
	// This member is required.
	DBClusterIdentifier *string

	// The name of the database engine to be used for this DB cluster. Valid Values:
	// neptune
	//
	// This member is required.
	Engine *string

	// A list of EC2 Availability Zones that instances in the DB cluster can be
	// created in.
	AvailabilityZones []string

	// The number of days for which automated backups are retained. You must specify a
	// minimum value of 1. Default: 1 Constraints:
	//   - Must be a value from 1 to 35
	BackupRetentionPeriod *int32

	// (Not supported by Neptune)
	CharacterSetName *string

	// If set to true , tags are copied to any snapshot of the DB cluster that is
	// created.
	CopyTagsToSnapshot *bool

	// The name of the DB cluster parameter group to associate with this DB cluster.
	// If this argument is omitted, the default is used. Constraints:
	//   - If supplied, must match the name of an existing DBClusterParameterGroup.
	DBClusterParameterGroupName *string

	// A DB subnet group to associate with this DB cluster. Constraints: Must match
	// the name of an existing DBSubnetGroup. Must not be default. Example:
	// mySubnetgroup
	DBSubnetGroupName *string

	// The name for your database of up to 64 alpha-numeric characters. If you do not
	// provide a name, Amazon Neptune will not create a database in the DB cluster you
	// are creating.
	DatabaseName *string

	// A value that indicates whether the DB cluster has deletion protection enabled.
	// The database can't be deleted when deletion protection is enabled. By default,
	// deletion protection is enabled.
	DeletionProtection *bool

	// The list of log types that need to be enabled for exporting to CloudWatch Logs.
	EnableCloudwatchLogsExports []string

	// If set to true , enables Amazon Identity and Access Management (IAM)
	// authentication for the entire DB cluster (this cannot be set at an instance
	// level). Default: false .
	EnableIAMDatabaseAuthentication *bool

	// The version number of the database engine to use for the new DB cluster.
	// Example: 1.0.2.1
	EngineVersion *string

	// The ID of the Neptune global database to which this new DB cluster should be
	// added.
	GlobalClusterIdentifier *string

	// The Amazon KMS key identifier for an encrypted DB cluster. The KMS key
	// identifier is the Amazon Resource Name (ARN) for the KMS encryption key. If you
	// are creating a DB cluster with the same Amazon account that owns the KMS
	// encryption key used to encrypt the new DB cluster, then you can use the KMS key
	// alias instead of the ARN for the KMS encryption key. If an encryption key is not
	// specified in KmsKeyId :
	//   - If ReplicationSourceIdentifier identifies an encrypted source, then Amazon
	//   Neptune will use the encryption key used to encrypt the source. Otherwise,
	//   Amazon Neptune will use your default encryption key.
	//   - If the StorageEncrypted parameter is true and ReplicationSourceIdentifier is
	//   not specified, then Amazon Neptune will use your default encryption key.
	// Amazon KMS creates the default encryption key for your Amazon account. Your
	// Amazon account has a different default encryption key for each Amazon Region. If
	// you create a Read Replica of an encrypted DB cluster in another Amazon Region,
	// you must set KmsKeyId to a KMS key ID that is valid in the destination Amazon
	// Region. This key is used to encrypt the Read Replica in that Amazon Region.
	KmsKeyId *string

	// Not supported by Neptune.
	MasterUserPassword *string

	// Not supported by Neptune.
	MasterUsername *string

	// (Not supported by Neptune)
	OptionGroupName *string

	// The port number on which the instances in the DB cluster accept connections.
	// Default: 8182
	Port *int32

	// This parameter is not currently supported.
	PreSignedUrl *string

	// The daily time range during which automated backups are created if automated
	// backups are enabled using the BackupRetentionPeriod parameter. The default is a
	// 30-minute window selected at random from an 8-hour block of time for each Amazon
	// Region. To see the time blocks available, see Adjusting the Preferred
	// Maintenance Window (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/AdjustingTheMaintenanceWindow.html)
	// in the Amazon Neptune User Guide. Constraints:
	//   - Must be in the format hh24:mi-hh24:mi .
	//   - Must be in Universal Coordinated Time (UTC).
	//   - Must not conflict with the preferred maintenance window.
	//   - Must be at least 30 minutes.
	PreferredBackupWindow *string

	// The weekly time range during which system maintenance can occur, in Universal
	// Coordinated Time (UTC). Format: ddd:hh24:mi-ddd:hh24:mi The default is a
	// 30-minute window selected at random from an 8-hour block of time for each Amazon
	// Region, occurring on a random day of the week. To see the time blocks available,
	// see Adjusting the Preferred Maintenance Window (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/AdjustingTheMaintenanceWindow.html)
	// in the Amazon Neptune User Guide. Valid Days: Mon, Tue, Wed, Thu, Fri, Sat, Sun.
	// Constraints: Minimum 30-minute window.
	PreferredMaintenanceWindow *string

	// The Amazon Resource Name (ARN) of the source DB instance or DB cluster if this
	// DB cluster is created as a Read Replica.
	ReplicationSourceIdentifier *string

	// Contains the scaling configuration of a Neptune Serverless DB cluster. For more
	// information, see Using Amazon Neptune Serverless (https://docs.aws.amazon.com/neptune/latest/userguide/neptune-serverless-using.html)
	// in the Amazon Neptune User Guide.
	ServerlessV2ScalingConfiguration *types.ServerlessV2ScalingConfiguration

	// The AWS region the resource is in. The presigned URL will be created with this
	// region, if the PresignURL member is empty set.
	SourceRegion *string

	// Specifies whether the DB cluster is encrypted.
	StorageEncrypted *bool

	// The tags to assign to the new DB cluster.
	Tags []types.Tag

	// A list of EC2 VPC security groups to associate with this DB cluster.
	VpcSecurityGroupIds []string

	// Used by the SDK's PresignURL autofill customization to specify the region the
	// of the client's request.
	destinationRegion *string

	noSmithyDocumentSerde
}

type CreateDBClusterOutput struct {

	// Contains the details of an Amazon Neptune DB cluster. This data type is used as
	// a response element in the DescribeDBClusters action.
	DBCluster *types.DBCluster

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreateDBClusterMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpCreateDBCluster{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpCreateDBCluster{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addCreateDBClusterPresignURLMiddleware(stack, options); err != nil {
		return err
	}
	if err = addCreateDBClusterResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpCreateDBClusterValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateDBCluster(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addendpointDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func copyCreateDBClusterInputForPresign(params interface{}) (interface{}, error) {
	input, ok := params.(*CreateDBClusterInput)
	if !ok {
		return nil, fmt.Errorf("expect *CreateDBClusterInput type, got %T", params)
	}
	cpy := *input
	return &cpy, nil
}
func getCreateDBClusterPreSignedUrl(params interface{}) (string, bool, error) {
	input, ok := params.(*CreateDBClusterInput)
	if !ok {
		return ``, false, fmt.Errorf("expect *CreateDBClusterInput type, got %T", params)
	}
	if input.PreSignedUrl == nil || len(*input.PreSignedUrl) == 0 {
		return ``, false, nil
	}
	return *input.PreSignedUrl, true, nil
}
func getCreateDBClusterSourceRegion(params interface{}) (string, bool, error) {
	input, ok := params.(*CreateDBClusterInput)
	if !ok {
		return ``, false, fmt.Errorf("expect *CreateDBClusterInput type, got %T", params)
	}
	if input.SourceRegion == nil || len(*input.SourceRegion) == 0 {
		return ``, false, nil
	}
	return *input.SourceRegion, true, nil
}
func setCreateDBClusterPreSignedUrl(params interface{}, value string) error {
	input, ok := params.(*CreateDBClusterInput)
	if !ok {
		return fmt.Errorf("expect *CreateDBClusterInput type, got %T", params)
	}
	input.PreSignedUrl = &value
	return nil
}
func setCreateDBClusterdestinationRegion(params interface{}, value string) error {
	input, ok := params.(*CreateDBClusterInput)
	if !ok {
		return fmt.Errorf("expect *CreateDBClusterInput type, got %T", params)
	}
	input.destinationRegion = &value
	return nil
}
func addCreateDBClusterPresignURLMiddleware(stack *middleware.Stack, options Options) error {
	return presignedurlcust.AddMiddleware(stack, presignedurlcust.Options{
		Accessor: presignedurlcust.ParameterAccessor{
			GetPresignedURL: getCreateDBClusterPreSignedUrl,

			GetSourceRegion: getCreateDBClusterSourceRegion,

			CopyInput: copyCreateDBClusterInputForPresign,

			SetDestinationRegion: setCreateDBClusterdestinationRegion,

			SetPresignedURL: setCreateDBClusterPreSignedUrl,
		},
		Presigner: &presignAutoFillCreateDBClusterClient{client: NewPresignClient(New(options))},
	})
}

type presignAutoFillCreateDBClusterClient struct {
	client *PresignClient
}

// PresignURL is a middleware accessor that satisfies URLPresigner interface.
func (c *presignAutoFillCreateDBClusterClient) PresignURL(ctx context.Context, srcRegion string, params interface{}) (*v4.PresignedHTTPRequest, error) {
	input, ok := params.(*CreateDBClusterInput)
	if !ok {
		return nil, fmt.Errorf("expect *CreateDBClusterInput type, got %T", params)
	}
	optFn := func(o *Options) {
		o.Region = srcRegion
		o.APIOptions = append(o.APIOptions, presignedurlcust.RemoveMiddleware)
	}
	presignOptFn := WithPresignClientFromClientOptions(optFn)
	return c.client.PresignCreateDBCluster(ctx, input, presignOptFn)
}

func newServiceMetadataMiddleware_opCreateDBCluster(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "rds",
		OperationName: "CreateDBCluster",
	}
}

// PresignCreateDBCluster is used to generate a presigned HTTP Request which
// contains presigned URL, signed headers and HTTP method used.
func (c *PresignClient) PresignCreateDBCluster(ctx context.Context, params *CreateDBClusterInput, optFns ...func(*PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	if params == nil {
		params = &CreateDBClusterInput{}
	}
	options := c.options.copy()
	for _, fn := range optFns {
		fn(&options)
	}
	clientOptFns := append(options.ClientOptions, withNopHTTPClientAPIOption)

	result, _, err := c.client.invokeOperation(ctx, "CreateDBCluster", params, clientOptFns,
		c.client.addOperationCreateDBClusterMiddlewares,
		presignConverter(options).convertToPresignMiddleware,
	)
	if err != nil {
		return nil, err
	}

	out := result.(*v4.PresignedHTTPRequest)
	return out, nil
}

type opCreateDBClusterResolveEndpointMiddleware struct {
	EndpointResolver EndpointResolverV2
	BuiltInResolver  builtInParameterResolver
}

func (*opCreateDBClusterResolveEndpointMiddleware) ID() string {
	return "ResolveEndpointV2"
}

func (m *opCreateDBClusterResolveEndpointMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	if awsmiddleware.GetRequiresLegacyEndpoints(ctx) {
		return next.HandleSerialize(ctx, in)
	}

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", in.Request)
	}

	if m.EndpointResolver == nil {
		return out, metadata, fmt.Errorf("expected endpoint resolver to not be nil")
	}

	params := EndpointParameters{}

	m.BuiltInResolver.ResolveBuiltIns(&params)

	var resolvedEndpoint smithyendpoints.Endpoint
	resolvedEndpoint, err = m.EndpointResolver.ResolveEndpoint(ctx, params)
	if err != nil {
		return out, metadata, fmt.Errorf("failed to resolve service endpoint, %w", err)
	}

	req.URL = &resolvedEndpoint.URI

	for k := range resolvedEndpoint.Headers {
		req.Header.Set(
			k,
			resolvedEndpoint.Headers.Get(k),
		)
	}

	authSchemes, err := internalauth.GetAuthenticationSchemes(&resolvedEndpoint.Properties)
	if err != nil {
		var nfe *internalauth.NoAuthenticationSchemesFoundError
		if errors.As(err, &nfe) {
			// if no auth scheme is found, default to sigv4
			signingName := "rds"
			signingRegion := m.BuiltInResolver.(*builtInResolver).Region
			ctx = awsmiddleware.SetSigningName(ctx, signingName)
			ctx = awsmiddleware.SetSigningRegion(ctx, signingRegion)

		}
		var ue *internalauth.UnSupportedAuthenticationSchemeSpecifiedError
		if errors.As(err, &ue) {
			return out, metadata, fmt.Errorf(
				"This operation requests signer version(s) %v but the client only supports %v",
				ue.UnsupportedSchemes,
				internalauth.SupportedSchemes,
			)
		}
	}

	for _, authScheme := range authSchemes {
		switch authScheme.(type) {
		case *internalauth.AuthenticationSchemeV4:
			v4Scheme, _ := authScheme.(*internalauth.AuthenticationSchemeV4)
			var signingName, signingRegion string
			if v4Scheme.SigningName == nil {
				signingName = "rds"
			} else {
				signingName = *v4Scheme.SigningName
			}
			if v4Scheme.SigningRegion == nil {
				signingRegion = m.BuiltInResolver.(*builtInResolver).Region
			} else {
				signingRegion = *v4Scheme.SigningRegion
			}
			if v4Scheme.DisableDoubleEncoding != nil {
				// The signer sets an equivalent value at client initialization time.
				// Setting this context value will cause the signer to extract it
				// and override the value set at client initialization time.
				ctx = internalauth.SetDisableDoubleEncoding(ctx, *v4Scheme.DisableDoubleEncoding)
			}
			ctx = awsmiddleware.SetSigningName(ctx, signingName)
			ctx = awsmiddleware.SetSigningRegion(ctx, signingRegion)
			break
		case *internalauth.AuthenticationSchemeV4A:
			v4aScheme, _ := authScheme.(*internalauth.AuthenticationSchemeV4A)
			if v4aScheme.SigningName == nil {
				v4aScheme.SigningName = aws.String("rds")
			}
			if v4aScheme.DisableDoubleEncoding != nil {
				// The signer sets an equivalent value at client initialization time.
				// Setting this context value will cause the signer to extract it
				// and override the value set at client initialization time.
				ctx = internalauth.SetDisableDoubleEncoding(ctx, *v4aScheme.DisableDoubleEncoding)
			}
			ctx = awsmiddleware.SetSigningName(ctx, *v4aScheme.SigningName)
			ctx = awsmiddleware.SetSigningRegion(ctx, v4aScheme.SigningRegionSet[0])
			break
		case *internalauth.AuthenticationSchemeNone:
			break
		}
	}

	return next.HandleSerialize(ctx, in)
}

func addCreateDBClusterResolveEndpointMiddleware(stack *middleware.Stack, options Options) error {
	return stack.Serialize.Insert(&opCreateDBClusterResolveEndpointMiddleware{
		EndpointResolver: options.EndpointResolverV2,
		BuiltInResolver: &builtInResolver{
			Region:       options.Region,
			UseDualStack: options.EndpointOptions.UseDualStackEndpoint,
			UseFIPS:      options.EndpointOptions.UseFIPSEndpoint,
			Endpoint:     options.BaseEndpoint,
		},
	}, "ResolveEndpoint", middleware.After)
}