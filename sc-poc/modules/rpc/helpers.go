package rpc

import (
	"fmt"
	"net/http"

	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

func processHTTPOptions(name string, opts *v1alpha1.HTTPOptions) v1alpha1.HTTPOptions {
	newOpts := v1alpha1.HTTPOptions{
		Method: http.MethodPost,
		URL:    fmt.Sprintf("/v1/api/complied-query/%s", name),
	}

	// Simply return if the provided options is nil
	if opts == nil {
		return newOpts
	}

	// Replace the options provided
	if opts.Method != "" {
		newOpts.Method = opts.Method
	}
	if opts.URL != "" {
		newOpts.URL = opts.URL
	}

	return newOpts
}
