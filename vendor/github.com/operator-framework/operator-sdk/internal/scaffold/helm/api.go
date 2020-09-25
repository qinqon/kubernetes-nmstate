// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helm

import (
	"fmt"
	"path/filepath"

	"github.com/prometheus/common/log"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/yaml"

	"github.com/operator-framework/operator-sdk/internal/genutil"
	"github.com/operator-framework/operator-sdk/internal/scaffold"
	"github.com/operator-framework/operator-sdk/internal/scaffold/input"
	"github.com/operator-framework/operator-sdk/pkg/helm/watches"
)

// TODO
// Consolidate scaffold func to be used by both "new" and "add api" commands.
func API(cfg input.Config, createOpts CreateChartOptions) error {

	r, chart, err := CreateChart(cfg.AbsProjectPath, createOpts)
	if err != nil {
		return fmt.Errorf("failed to create helm chart: %v", err)
	}

	valuesPath := filepath.Join("<project_dir>", HelmChartsDir, chart.Name(), "values.yaml")

	rawValues, err := yaml.Marshal(chart.Values)
	if err != nil {
		return fmt.Errorf("failed to get raw chart values: %v", err)
	}
	crSpec := fmt.Sprintf("# Default values copied from %s\n\n%s", valuesPath, rawValues)

	// update watch.yaml for the given resource.
	watchesFile := filepath.Join(cfg.AbsProjectPath, watches.WatchesFile)
	if err := watches.UpdateForResource(watchesFile, r, chart.Name()); err != nil {
		return fmt.Errorf("failed to update watches.yaml: %w", err)
	}

	s := &scaffold.Scaffold{}
	err = s.Execute(&cfg,
		&scaffold.CR{
			Resource: r,
			Spec:     crSpec,
		},
	)
	if err != nil {
		log.Fatalf("API scaffold failed: %v", err)
	}
	if err = genutil.GenerateCRDNonGo("", *r, createOpts.CRDVersion); err != nil {
		return err
	}

	roleScaffold := DefaultRoleScaffold

	if k8sCfg, err := config.GetConfig(); err != nil {
		log.Warnf("Using default RBAC rules: failed to get Kubernetes config: %s", err)
	} else if dc, err := discovery.NewDiscoveryClientForConfig(k8sCfg); err != nil {
		log.Warnf("Using default RBAC rules: failed to create Kubernetes discovery client: %s", err)
	} else {
		roleScaffold = GenerateRoleScaffold(dc, chart)
	}

	if err = scaffold.MergeRoleForResource(r, cfg.AbsProjectPath, roleScaffold); err != nil {
		return fmt.Errorf("failed to merge rules in the RBAC manifest for resource (%v, %v): %v",
			r.APIVersion, r.Kind, err)
	}

	return nil
}
