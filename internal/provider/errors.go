package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ErrUnknownProviderAttribute(attributeName, environmentName string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path.Root(attributeName),
		fmt.Sprintf("Unknown %s", attributeName),
		fmt.Sprintf("The provider cannot create the Octopus Deploy API client as %s resolved to an unknown configuration value. "+
			"Either target apply the source of the value first, set the value statically in the configuration, or use the %s environment variable.",
			attributeName,
			environmentName,
		),
	)
}

func ErrMissingProviderAttribute(attributeName, environmentName string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path.Root(attributeName),
		fmt.Sprintf("Missing %s", attributeName),
		fmt.Sprintf("The provider cannot create the Octopus Deploy API client as there is a missing or empty value for the Octopus Deploy %s. "+
			"Set %s in the configuration or use the %s environment variable.",
			attributeName,
			attributeName,
			environmentName,
		),
	)
}

func ErrMissingAttribute(attributeName string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path.Root(attributeName),
		fmt.Sprintf("Missing %s", attributeName),
		fmt.Sprintf("Must provide %s.", attributeName),
	)
}

func ErrUnexpectedDataConfigureType(input any) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Unexpected Data Source Configure Type",
		fmt.Sprintf(
			"Expected *api.Client, got: %T. Please report this issue to the provider developers.",
			input,
		),
	)
}
