// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package cmd

import (
	"fmt"
	"strings"

	"github.com/Azure/aks-engine/pkg/api"
	"github.com/Azure/aks-engine/pkg/helpers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	orchestratorsName             = "orchestrators"
	orchestratorsShortDescription = "Display info about supported orchestrators"
	orchestratorsLongDescription  = "Display supported versions and upgrade versions for each orchestrator"
)

type orchestratorsCmd struct {
	// user input
	orchestrator string
	version      string
	cloudType    string
	windows      bool
}

func newOrchestratorsCmd() *cobra.Command {
	oc := orchestratorsCmd{}

	command := &cobra.Command{
		Use:   orchestratorsName,
		Short: orchestratorsShortDescription,
		Long:  orchestratorsLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := oc.validate(cmd, args); err != nil {
				log.Fatalf(fmt.Sprintf("error validating orchestratorsCmd: %s", err.Error()))
			}
			return oc.run(cmd, args)
		},
	}

	f := command.Flags()
	f.StringVar(&oc.orchestrator, "orchestrator", "", "orchestrator name (optional) ")
	f.StringVar(&oc.version, "version", "", "orchestrator version (optional)")
	f.StringVar(&oc.cloudType, "cloud-type", "azure", "cloud type. the value should be either azure or azurestack (optional)")
	f.BoolVar(&oc.windows, "windows", false, "orchestrator platform (optional, applies to Kubernetes only)")

	return command
}

func (oc *orchestratorsCmd) validate(cmd *cobra.Command, args []string) error {

	if !(strings.EqualFold(oc.cloudType, api.AzureCloudType) || strings.EqualFold(oc.cloudType, api.AzureStackCloudType)) {
		return errors.Errorf("--cloud-type value (%s) is invalid. The value should be either azure or azurestack", oc.cloudType)
	}

	return nil
}

func (oc *orchestratorsCmd) run(cmd *cobra.Command, args []string) error {
	orchs, err := api.GetOrchestratorVersionProfileListVLabs(oc.orchestrator, oc.version, oc.windows, oc.cloudType)
	if err != nil {
		return err
	}

	data, err := helpers.JSONMarshalIndent(orchs, "", "  ", false)
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}
