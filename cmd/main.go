/*
Copyright 2023 Clay.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	goflag "flag"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	klog "k8s.io/klog/v2"

	etcdclustercontroller "etcd-operator/cmd/etcdcluster-controller"
	etcdinspectioncontroller "etcd-operator/cmd/etcdinspection-controller"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var cmd = &cobra.Command{
		Use:              "etcd-controller",
		Short:            "run etcdcluster/etcdinspection controller",
		PersistentPreRun: func(c *cobra.Command, args []string) {},
	}

	flags := cmd.PersistentFlags()
	out := cmd.OutOrStdout()
	cmd.AddCommand(
		etcdclustercontroller.NewEtcdClusterControllerCommand(out),
		etcdinspectioncontroller.NewEtcdInspectionControllerCommand(out),
	)

	klog.InitFlags(nil)
	defer klog.Flush()

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	err := flags.Parse(os.Args[1:])
	if err != nil {
		klog.Errorf("failed to parse args, err is %v", err)
	}

	if err = cmd.Execute(); err != nil {
		klog.Errorf("failed to execute controller, err: %v", err)
		os.Exit(1)
	}
}
