// Code generated by smithy-go-codegen DO NOT EDIT.

package lambda

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	internalauth "github.com/aws/aws-sdk-go-v2/internal/auth"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates a Lambda function. To create a function, you need a deployment package (https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-package.html)
// and an execution role (https://docs.aws.amazon.com/lambda/latest/dg/intro-permission-model.html#lambda-intro-execution-role)
// . The deployment package is a .zip file archive or container image that contains
// your function code. The execution role grants the function permission to use
// Amazon Web Services, such as Amazon CloudWatch Logs for log streaming and X-Ray
// for request tracing. If the deployment package is a container image (https://docs.aws.amazon.com/lambda/latest/dg/lambda-images.html)
// , then you set the package type to Image . For a container image, the code
// property must include the URI of a container image in the Amazon ECR registry.
// You do not need to specify the handler and runtime properties. If the deployment
// package is a .zip file archive (https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-package.html#gettingstarted-package-zip)
// , then you set the package type to Zip . For a .zip file archive, the code
// property specifies the location of the .zip file. You must also specify the
// handler and runtime properties. The code in the deployment package must be
// compatible with the target instruction set architecture of the function ( x86-64
// or arm64 ). If you do not specify the architecture, then the default value is
// x86-64 . When you create a function, Lambda provisions an instance of the
// function and its supporting resources. If your function connects to a VPC, this
// process can take a minute or so. During this time, you can't invoke or modify
// the function. The State , StateReason , and StateReasonCode fields in the
// response from GetFunctionConfiguration indicate when the function is ready to
// invoke. For more information, see Lambda function states (https://docs.aws.amazon.com/lambda/latest/dg/functions-states.html)
// . A function has an unpublished version, and can have published versions and
// aliases. The unpublished version changes when you update your function's code
// and configuration. A published version is a snapshot of your function code and
// configuration that can't be changed. An alias is a named resource that maps to a
// version, and can be changed to map to a different version. Use the Publish
// parameter to create version 1 of your function from its initial configuration.
// The other parameters let you configure version-specific and function-level
// settings. You can modify version-specific settings later with
// UpdateFunctionConfiguration . Function-level settings apply to both the
// unpublished and published versions of the function, and include tags (
// TagResource ) and per-function concurrency limits ( PutFunctionConcurrency ).
// You can use code signing if your deployment package is a .zip file archive. To
// enable code signing for this function, specify the ARN of a code-signing
// configuration. When a user attempts to deploy a code package with
// UpdateFunctionCode , Lambda checks that the code package has a valid signature
// from a trusted publisher. The code-signing configuration includes set of signing
// profiles, which define the trusted publishers for this function. If another
// Amazon Web Services account or an Amazon Web Service invokes your function, use
// AddPermission to grant permission by creating a resource-based Identity and
// Access Management (IAM) policy. You can grant permissions at the function level,
// on a version, or on an alias. To invoke your function directly, use Invoke . To
// invoke your function in response to events in other Amazon Web Services, create
// an event source mapping ( CreateEventSourceMapping ), or configure a function
// trigger in the other service. For more information, see Invoking Lambda
// functions (https://docs.aws.amazon.com/lambda/latest/dg/lambda-invocation.html) .
func (c *Client) CreateFunction(ctx context.Context, params *CreateFunctionInput, optFns ...func(*Options)) (*CreateFunctionOutput, error) {
	if params == nil {
		params = &CreateFunctionInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateFunction", params, optFns, c.addOperationCreateFunctionMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateFunctionOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateFunctionInput struct {

	// The code for the function.
	//
	// This member is required.
	Code *types.FunctionCode

	// The name of the Lambda function. Name formats
	//   - Function name – my-function .
	//   - Function ARN – arn:aws:lambda:us-west-2:123456789012:function:my-function .
	//   - Partial ARN – 123456789012:function:my-function .
	// The length constraint applies only to the full ARN. If you specify only the
	// function name, it is limited to 64 characters in length.
	//
	// This member is required.
	FunctionName *string

	// The Amazon Resource Name (ARN) of the function's execution role.
	//
	// This member is required.
	Role *string

	// The instruction set architecture that the function supports. Enter a string
	// array with one of the valid values (arm64 or x86_64). The default value is
	// x86_64 .
	Architectures []types.Architecture

	// To enable code signing for this function, specify the ARN of a code-signing
	// configuration. A code-signing configuration includes a set of signing profiles,
	// which define the trusted publishers for this function.
	CodeSigningConfigArn *string

	// A dead-letter queue configuration that specifies the queue or topic where
	// Lambda sends asynchronous events when they fail processing. For more
	// information, see Dead-letter queues (https://docs.aws.amazon.com/lambda/latest/dg/invocation-async.html#invocation-dlq)
	// .
	DeadLetterConfig *types.DeadLetterConfig

	// A description of the function.
	Description *string

	// Environment variables that are accessible from function code during execution.
	Environment *types.Environment

	// The size of the function's /tmp directory in MB. The default value is 512, but
	// can be any whole number between 512 and 10,240 MB.
	EphemeralStorage *types.EphemeralStorage

	// Connection settings for an Amazon EFS file system.
	FileSystemConfigs []types.FileSystemConfig

	// The name of the method within your code that Lambda calls to run your function.
	// Handler is required if the deployment package is a .zip file archive. The format
	// includes the file name. It can also include namespaces and other qualifiers,
	// depending on the runtime. For more information, see Lambda programming model (https://docs.aws.amazon.com/lambda/latest/dg/foundation-progmodel.html)
	// .
	Handler *string

	// Container image configuration values (https://docs.aws.amazon.com/lambda/latest/dg/configuration-images.html#configuration-images-settings)
	// that override the values in the container image Dockerfile.
	ImageConfig *types.ImageConfig

	// The ARN of the Key Management Service (KMS) customer managed key that's used to
	// encrypt your function's environment variables (https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-encryption)
	// . When Lambda SnapStart (https://docs.aws.amazon.com/lambda/latest/dg/snapstart-security.html)
	// is activated, Lambda also uses this key is to encrypt your function's snapshot.
	// If you deploy your function using a container image, Lambda also uses this key
	// to encrypt your function when it's deployed. Note that this is not the same key
	// that's used to protect your container image in the Amazon Elastic Container
	// Registry (Amazon ECR). If you don't provide a customer managed key, Lambda uses
	// a default service key.
	KMSKeyArn *string

	// A list of function layers (https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html)
	// to add to the function's execution environment. Specify each layer by its ARN,
	// including the version.
	Layers []string

	// The amount of memory available to the function (https://docs.aws.amazon.com/lambda/latest/dg/configuration-function-common.html#configuration-memory-console)
	// at runtime. Increasing the function memory also increases its CPU allocation.
	// The default value is 128 MB. The value can be any multiple of 1 MB.
	MemorySize *int32

	// The type of deployment package. Set to Image for container image and set to Zip
	// for .zip file archive.
	PackageType types.PackageType

	// Set to true to publish the first version of the function during creation.
	Publish bool

	// The identifier of the function's runtime (https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html)
	// . Runtime is required if the deployment package is a .zip file archive. The
	// following list includes deprecated runtimes. For more information, see Runtime
	// deprecation policy (https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html#runtime-support-policy)
	// .
	Runtime types.Runtime

	// The function's SnapStart (https://docs.aws.amazon.com/lambda/latest/dg/snapstart.html)
	// setting.
	SnapStart *types.SnapStart

	// A list of tags (https://docs.aws.amazon.com/lambda/latest/dg/tagging.html) to
	// apply to the function.
	Tags map[string]string

	// The amount of time (in seconds) that Lambda allows a function to run before
	// stopping it. The default is 3 seconds. The maximum allowed value is 900 seconds.
	// For more information, see Lambda execution environment (https://docs.aws.amazon.com/lambda/latest/dg/runtimes-context.html)
	// .
	Timeout *int32

	// Set Mode to Active to sample and trace a subset of incoming requests with X-Ray (https://docs.aws.amazon.com/lambda/latest/dg/services-xray.html)
	// .
	TracingConfig *types.TracingConfig

	// For network connectivity to Amazon Web Services resources in a VPC, specify a
	// list of security groups and subnets in the VPC. When you connect a function to a
	// VPC, it can access resources and the internet only through that VPC. For more
	// information, see Configuring a Lambda function to access resources in a VPC (https://docs.aws.amazon.com/lambda/latest/dg/configuration-vpc.html)
	// .
	VpcConfig *types.VpcConfig

	noSmithyDocumentSerde
}

// Details about a function's configuration.
type CreateFunctionOutput struct {

	// The instruction set architecture that the function supports. Architecture is a
	// string array with one of the valid values. The default architecture value is
	// x86_64 .
	Architectures []types.Architecture

	// The SHA256 hash of the function's deployment package.
	CodeSha256 *string

	// The size of the function's deployment package, in bytes.
	CodeSize int64

	// The function's dead letter queue.
	DeadLetterConfig *types.DeadLetterConfig

	// The function's description.
	Description *string

	// The function's environment variables (https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html)
	// . Omitted from CloudTrail logs.
	Environment *types.EnvironmentResponse

	// The size of the function’s /tmp directory in MB. The default value is 512, but
	// it can be any whole number between 512 and 10,240 MB.
	EphemeralStorage *types.EphemeralStorage

	// Connection settings for an Amazon EFS file system (https://docs.aws.amazon.com/lambda/latest/dg/configuration-filesystem.html)
	// .
	FileSystemConfigs []types.FileSystemConfig

	// The function's Amazon Resource Name (ARN).
	FunctionArn *string

	// The name of the function.
	FunctionName *string

	// The function that Lambda calls to begin running your function.
	Handler *string

	// The function's image configuration values.
	ImageConfigResponse *types.ImageConfigResponse

	// The KMS key that's used to encrypt the function's environment variables (https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-encryption)
	// . When Lambda SnapStart (https://docs.aws.amazon.com/lambda/latest/dg/snapstart-security.html)
	// is activated, this key is also used to encrypt the function's snapshot. This key
	// is returned only if you've configured a customer managed key.
	KMSKeyArn *string

	// The date and time that the function was last updated, in ISO-8601 format (https://www.w3.org/TR/NOTE-datetime)
	// (YYYY-MM-DDThh:mm:ss.sTZD).
	LastModified *string

	// The status of the last update that was performed on the function. This is first
	// set to Successful after function creation completes.
	LastUpdateStatus types.LastUpdateStatus

	// The reason for the last update that was performed on the function.
	LastUpdateStatusReason *string

	// The reason code for the last update that was performed on the function.
	LastUpdateStatusReasonCode types.LastUpdateStatusReasonCode

	// The function's layers (https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html)
	// .
	Layers []types.Layer

	// For Lambda@Edge functions, the ARN of the main function.
	MasterArn *string

	// The amount of memory available to the function at runtime.
	MemorySize *int32

	// The type of deployment package. Set to Image for container image and set Zip
	// for .zip file archive.
	PackageType types.PackageType

	// The latest updated revision of the function or alias.
	RevisionId *string

	// The function's execution role.
	Role *string

	// The identifier of the function's runtime (https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html)
	// . Runtime is required if the deployment package is a .zip file archive. The
	// following list includes deprecated runtimes. For more information, see Runtime
	// deprecation policy (https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html#runtime-support-policy)
	// .
	Runtime types.Runtime

	// The ARN of the runtime and any errors that occured.
	RuntimeVersionConfig *types.RuntimeVersionConfig

	// The ARN of the signing job.
	SigningJobArn *string

	// The ARN of the signing profile version.
	SigningProfileVersionArn *string

	// Set ApplyOn to PublishedVersions to create a snapshot of the initialized
	// execution environment when you publish a function version. For more information,
	// see Improving startup performance with Lambda SnapStart (https://docs.aws.amazon.com/lambda/latest/dg/snapstart.html)
	// .
	SnapStart *types.SnapStartResponse

	// The current state of the function. When the state is Inactive , you can
	// reactivate the function by invoking it.
	State types.State

	// The reason for the function's current state.
	StateReason *string

	// The reason code for the function's current state. When the code is Creating ,
	// you can't invoke or modify the function.
	StateReasonCode types.StateReasonCode

	// The amount of time in seconds that Lambda allows a function to run before
	// stopping it.
	Timeout *int32

	// The function's X-Ray tracing configuration.
	TracingConfig *types.TracingConfigResponse

	// The version of the Lambda function.
	Version *string

	// The function's networking configuration.
	VpcConfig *types.VpcConfigResponse

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreateFunctionMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsRestjson1_serializeOpCreateFunction{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpCreateFunction{}, middleware.After)
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
	if err = addCreateFunctionResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpCreateFunctionValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateFunction(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opCreateFunction(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "lambda",
		OperationName: "CreateFunction",
	}
}

type opCreateFunctionResolveEndpointMiddleware struct {
	EndpointResolver EndpointResolverV2
	BuiltInResolver  builtInParameterResolver
}

func (*opCreateFunctionResolveEndpointMiddleware) ID() string {
	return "ResolveEndpointV2"
}

func (m *opCreateFunctionResolveEndpointMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
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
			signingName := "lambda"
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
				signingName = "lambda"
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
				v4aScheme.SigningName = aws.String("lambda")
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

func addCreateFunctionResolveEndpointMiddleware(stack *middleware.Stack, options Options) error {
	return stack.Serialize.Insert(&opCreateFunctionResolveEndpointMiddleware{
		EndpointResolver: options.EndpointResolverV2,
		BuiltInResolver: &builtInResolver{
			Region:       options.Region,
			UseDualStack: options.EndpointOptions.UseDualStackEndpoint,
			UseFIPS:      options.EndpointOptions.UseFIPSEndpoint,
			Endpoint:     options.BaseEndpoint,
		},
	}, "ResolveEndpoint", middleware.After)
}