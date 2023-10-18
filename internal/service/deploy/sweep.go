// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package deploy

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv1"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_codedeploy_app", &resource.Sweeper{
		Name: "aws_codedeploy_app",
		F:    sweepApps,
	})
}

func sweepApps(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)

	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.DeployConn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	input := &codedeploy.ListApplicationsInput{}

	err = conn.ListApplicationsPagesWithContext(ctx, input, func(page *codedeploy.ListApplicationsOutput, lastPage bool) bool {
		for _, app := range page.Applications {
			if app == nil {
				continue
			}

			appName := aws.StringValue(app)
			r := ResourceApp()
			d := r.Data(nil)
			d.SetId(fmt.Sprintf("%s:%s", "xxxx", appName))
			d.Set("name", appName)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing CodeDeploy Applications for %s: %w", region, err))
	}

	if err = sweep.SweepOrchestrator(ctx, sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping CodeDeploy Applications for %s: %w", region, err))
	}

	if awsv1.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping CodeDeploy Applications sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}
