package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func expandStringList(ctx context.Context, in types.List) ([]string, diag.Diagnostics) {
	count := len(in.Elements())
	vals := make([]types.String, 0, count)
	if count < 1 {
		return []string{}, nil
	}

	diags := in.ElementsAs(ctx, &vals, false)
	if diags.HasError() {
		return []string{}, diags
	}

	out := make([]string, 0, len(vals))
	for _, val := range vals {
		out = append(out, val.ValueString())
	}

	return out, diags
}

func flattenStringList(ctx context.Context, in []string) (types.List, diag.Diagnostics) {
	if in == nil {
		in = []string{}
	}

	out, diags := types.ListValueFrom(ctx, types.StringType, in)
	return out, diags
}
