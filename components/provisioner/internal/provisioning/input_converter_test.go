package provisioning

import (
	"testing"

	"github.com/kyma-incubator/compass/components/provisioner/internal/util"
	"github.com/kyma-incubator/compass/components/provisioner/internal/uuid"

	"github.com/kyma-incubator/compass/components/provisioner/internal/installation/release"
	realeaseMocks "github.com/kyma-incubator/compass/components/provisioner/internal/installation/release/mocks"

	"github.com/kyma-incubator/compass/components/provisioner/internal/persistence/dberrors"

	"github.com/kyma-incubator/compass/components/provisioner/internal/uuid/mocks"
	"github.com/stretchr/testify/require"

	"github.com/kyma-incubator/compass/components/provisioner/internal/model"

	"github.com/kyma-incubator/compass/components/provisioner/pkg/gqlschema"
	"github.com/stretchr/testify/assert"
)

const (
	kymaVersion                   = "1.5"
	clusterEssentialsComponent    = "cluster-essentials"
	coreComponent                 = "core"
	rafterComponent               = "rafter"
	applicationConnectorComponent = "application-connector"

	rafterSourceURL = "github.com/kyma-project/kyma.git//resources/rafter"

	gardenerProject = "gardener-project"
)

func Test_ProvisioningInputToCluster(t *testing.T) {

	readSession := &realeaseMocks.Repository{}
	readSession.On("GetReleaseByVersion", kymaVersion).Return(fixKymaRelease(), nil)

	createGQLRuntimeInputGCP := func(zone *string) gqlschema.ProvisionRuntimeInput {
		return gqlschema.ProvisionRuntimeInput{
			RuntimeInput: &gqlschema.RuntimeInput{
				Name:        "runtimeName",
				Description: nil,
				Labels:      &gqlschema.Labels{},
			},
			ClusterConfig: &gqlschema.ClusterConfigInput{
				GcpConfig: &gqlschema.GCPConfigInput{
					Name:              "Something",
					ProjectName:       "Project",
					NumberOfNodes:     3,
					BootDiskSizeGb:    256,
					MachineType:       "n1-standard-1",
					Region:            "region",
					Zone:              zone,
					KubernetesVersion: "version",
				},
			},
			Credentials: &gqlschema.CredentialsInput{
				SecretName: "secretName",
			},
			KymaConfig: fixKymaGraphQLConfigInput(),
		}
	}

	createExpectedRuntimeInputGCP := func(zone string) model.Cluster {
		return model.Cluster{
			ID: "runtimeID",
			ClusterConfig: model.GCPConfig{
				ID:                "id",
				Name:              "Something",
				ProjectName:       "Project",
				NumberOfNodes:     3,
				BootDiskSizeGB:    256,
				MachineType:       "n1-standard-1",
				Region:            "region",
				Zone:              zone,
				KubernetesVersion: "version",
				ClusterID:         "runtimeID",
			},
			Kubeconfig:            nil,
			KymaConfig:            fixKymaConfig(),
			CredentialsSecretName: "",
			Tenant:                tenant,
			SubAccountId:          util.StringPtr(subAccountId),
		}
	}

	gcpGardenerProvider := &gqlschema.GCPProviderConfigInput{Zones: []string{"fix-gcp-zone-1", "fix-gcp-zone-2"}}

	gardenerGCPGQLInput := gqlschema.ProvisionRuntimeInput{
		RuntimeInput: &gqlschema.RuntimeInput{
			Name:        "runtimeName",
			Description: nil,
			Labels:      &gqlschema.Labels{},
		},
		ClusterConfig: &gqlschema.ClusterConfigInput{
			GardenerConfig: &gqlschema.GardenerConfigInput{
				KubernetesVersion: "version",
				VolumeSizeGb:      1024,
				MachineType:       "n1-standard-1",
				Region:            "region",
				Provider:          "GCP",
				Seed:              util.StringPtr("gcp-eu1"),
				TargetSecret:      "secret",
				DiskType:          "ssd",
				WorkerCidr:        "cidr",
				AutoScalerMin:     1,
				AutoScalerMax:     5,
				MaxSurge:          1,
				MaxUnavailable:    2,
				ProviderSpecificConfig: &gqlschema.ProviderSpecificInput{
					GcpConfig: gcpGardenerProvider,
				},
			},
		},
		Credentials: &gqlschema.CredentialsInput{
			SecretName: "secretName",
		},
		KymaConfig: fixKymaGraphQLConfigInput(),
	}

	expectedGCPProviderCfg, err := model.NewGCPGardenerConfig(gcpGardenerProvider)
	require.NoError(t, err)

	expectedGardenerGCPRuntimeConfig := model.Cluster{
		ID: "runtimeID",
		ClusterConfig: model.GardenerConfig{
			ID:                     "id",
			Name:                   "verylon",
			ProjectName:            gardenerProject,
			MachineType:            "n1-standard-1",
			Region:                 "region",
			KubernetesVersion:      "version",
			VolumeSizeGB:           1024,
			DiskType:               "ssd",
			Provider:               "GCP",
			Seed:                   "gcp-eu1",
			TargetSecret:           "secret",
			WorkerCidr:             "cidr",
			AutoScalerMin:          1,
			AutoScalerMax:          5,
			MaxSurge:               1,
			MaxUnavailable:         2,
			ClusterID:              "runtimeID",
			GardenerProviderConfig: expectedGCPProviderCfg,
		},
		Kubeconfig:            nil,
		KymaConfig:            fixKymaConfig(),
		CredentialsSecretName: "",
		Tenant:                tenant,
		SubAccountId:          util.StringPtr(subAccountId),
	}

	createGQLRuntimeInputAzure := func(zones []string) gqlschema.ProvisionRuntimeInput {
		return gqlschema.ProvisionRuntimeInput{
			RuntimeInput: &gqlschema.RuntimeInput{
				Name:        "runtimeName",
				Description: nil,
				Labels:      &gqlschema.Labels{},
			},
			ClusterConfig: &gqlschema.ClusterConfigInput{
				GardenerConfig: &gqlschema.GardenerConfigInput{
					KubernetesVersion: "version",
					VolumeSizeGb:      1024,
					MachineType:       "n1-standard-1",
					Region:            "region",
					Provider:          "Azure",
					TargetSecret:      "secret",
					DiskType:          "ssd",
					WorkerCidr:        "cidr",
					AutoScalerMin:     1,
					AutoScalerMax:     5,
					MaxSurge:          1,
					MaxUnavailable:    2,
					ProviderSpecificConfig: &gqlschema.ProviderSpecificInput{
						AzureConfig: &gqlschema.AzureProviderConfigInput{
							VnetCidr: "cidr",
							Zones:    zones,
						},
					},
				},
			},
			Credentials: &gqlschema.CredentialsInput{
				SecretName: "secretName",
			},
			KymaConfig: fixKymaGraphQLConfigInput(),
		}
	}

	expectedGardenerAzureRuntimeConfig := func(zones []string) model.Cluster {

		expectedAzureProviderCfg, err := model.NewAzureGardenerConfig(&gqlschema.AzureProviderConfigInput{VnetCidr: "cidr", Zones: zones})
		require.NoError(t, err)

		return model.Cluster{
			ID: "runtimeID",
			ClusterConfig: model.GardenerConfig{
				ID:                     "id",
				Name:                   "verylon",
				ProjectName:            gardenerProject,
				MachineType:            "n1-standard-1",
				Region:                 "region",
				KubernetesVersion:      "version",
				VolumeSizeGB:           1024,
				DiskType:               "ssd",
				Provider:               "Azure",
				Seed:                   "",
				TargetSecret:           "secret",
				WorkerCidr:             "cidr",
				AutoScalerMin:          1,
				AutoScalerMax:          5,
				MaxSurge:               1,
				MaxUnavailable:         2,
				ClusterID:              "runtimeID",
				GardenerProviderConfig: expectedAzureProviderCfg,
			},
			Kubeconfig:            nil,
			KymaConfig:            fixKymaConfig(),
			CredentialsSecretName: "",
			Tenant:                tenant,
			SubAccountId:          util.StringPtr(subAccountId),
		}
	}

	awsGardenerProvider := &gqlschema.AWSProviderConfigInput{
		Zone:         "zone",
		InternalCidr: "cidr",
		VpcCidr:      "cidr",
		PublicCidr:   "cidr",
	}

	gardenerAWSGQLInput := gqlschema.ProvisionRuntimeInput{
		RuntimeInput: &gqlschema.RuntimeInput{
			Name:        "runtimeName",
			Description: nil,
			Labels:      &gqlschema.Labels{},
		},
		ClusterConfig: &gqlschema.ClusterConfigInput{
			GardenerConfig: &gqlschema.GardenerConfigInput{
				KubernetesVersion: "version",
				VolumeSizeGb:      1024,
				MachineType:       "n1-standard-1",
				Region:            "region",
				Provider:          "AWS",
				Seed:              util.StringPtr("aws-eu1"),
				TargetSecret:      "secret",
				DiskType:          "ssd",
				WorkerCidr:        "cidr",
				AutoScalerMin:     1,
				AutoScalerMax:     5,
				MaxSurge:          1,
				MaxUnavailable:    2,
				ProviderSpecificConfig: &gqlschema.ProviderSpecificInput{
					AwsConfig: awsGardenerProvider,
				},
			},
		},
		Credentials: &gqlschema.CredentialsInput{
			SecretName: "secretName",
		},
		KymaConfig: fixKymaGraphQLConfigInput(),
	}

	expectedAWSProviderCfg, err := model.NewAWSGardenerConfig(awsGardenerProvider)
	require.NoError(t, err)

	expectedGardenerAWSRuntimeConfig := model.Cluster{
		ID: "runtimeID",
		ClusterConfig: model.GardenerConfig{
			ID:                     "id",
			Name:                   "verylon",
			ProjectName:            gardenerProject,
			MachineType:            "n1-standard-1",
			Region:                 "region",
			KubernetesVersion:      "version",
			VolumeSizeGB:           1024,
			DiskType:               "ssd",
			Provider:               "AWS",
			Seed:                   "aws-eu1",
			TargetSecret:           "secret",
			WorkerCidr:             "cidr",
			AutoScalerMin:          1,
			AutoScalerMax:          5,
			MaxSurge:               1,
			MaxUnavailable:         2,
			ClusterID:              "runtimeID",
			GardenerProviderConfig: expectedAWSProviderCfg,
		},
		Kubeconfig:            nil,
		KymaConfig:            fixKymaConfig(),
		CredentialsSecretName: "",
		Tenant:                tenant,
		SubAccountId:          util.StringPtr(subAccountId),
	}

	zone := "zone"
	gardenerZones := []string{"fix-az-zone-1", "fix-az-zone-2"}

	configurations := []struct {
		input       gqlschema.ProvisionRuntimeInput
		expected    model.Cluster
		description string
	}{
		{
			input:       createGQLRuntimeInputGCP(&zone),
			expected:    createExpectedRuntimeInputGCP(zone),
			description: "Should create proper runtime config struct with GCP input",
		},
		{
			input:       createGQLRuntimeInputGCP(nil),
			expected:    createExpectedRuntimeInputGCP(""),
			description: "Should create proper runtime config struct with GCP input (empty zone)",
		},
		{
			input:       gardenerGCPGQLInput,
			expected:    expectedGardenerGCPRuntimeConfig,
			description: "Should create proper runtime config struct with Gardener input for GCP provider",
		},
		{
			input:       createGQLRuntimeInputAzure(nil),
			expected:    expectedGardenerAzureRuntimeConfig(nil),
			description: "Should create proper runtime config struct with Gardener input for Azure provider",
		},
		{
			input:       createGQLRuntimeInputAzure(gardenerZones),
			expected:    expectedGardenerAzureRuntimeConfig(gardenerZones),
			description: "Should create proper runtime config struct with Gardener input for Azure provider with zones passed",
		},
		{
			input:       gardenerAWSGQLInput,
			expected:    expectedGardenerAWSRuntimeConfig,
			description: "Should create proper runtime config struct with Gardener input for AWS provider",
		},
	}

	for _, testCase := range configurations {
		t.Run(testCase.description, func(t *testing.T) {
			//given
			uuidGeneratorMock := &mocks.UUIDGenerator{}
			uuidGeneratorMock.On("New").Return("id").Times(6)
			uuidGeneratorMock.On("New").Return("very-Long-ID-That-Has-More-Than-Fourteen-Characters-And-Even-Some-Hyphens")

			inputConverter := NewInputConverter(uuidGeneratorMock, readSession, gardenerProject)

			//when
			runtimeConfig, err := inputConverter.ProvisioningInputToCluster("runtimeID", testCase.input, tenant, subAccountId)

			//then
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, runtimeConfig)
			uuidGeneratorMock.AssertExpectations(t)
		})
	}
}

func TestConverter_ProvisioningInputToCluster_Error(t *testing.T) {

	t.Run("should return error when failed to get kyma release", func(t *testing.T) {
		// given
		uuidGeneratorMock := &mocks.UUIDGenerator{}
		readSession := &realeaseMocks.Repository{}
		readSession.On("GetReleaseByVersion", kymaVersion).Return(model.Release{}, dberrors.NotFound("error"))

		input := gqlschema.ProvisionRuntimeInput{
			ClusterConfig: &gqlschema.ClusterConfigInput{
				GcpConfig: &gqlschema.GCPConfigInput{},
			},
			Credentials: &gqlschema.CredentialsInput{
				SecretName: "secretName",
			},
			KymaConfig: &gqlschema.KymaConfigInput{
				Version: kymaVersion,
			},
		}

		inputConverter := NewInputConverter(uuidGeneratorMock, readSession, gardenerProject)

		//when
		_, err := inputConverter.ProvisioningInputToCluster("runtimeID", input, tenant, subAccountId)

		//then
		require.Error(t, err)
		uuidGeneratorMock.AssertExpectations(t)
	})

	t.Run("should return error when no provider input provided", func(t *testing.T) {
		// given
		input := gqlschema.ProvisionRuntimeInput{
			ClusterConfig: &gqlschema.ClusterConfigInput{},
		}

		inputConverter := NewInputConverter(nil, nil, gardenerProject)

		//when
		_, err := inputConverter.ProvisioningInputToCluster("runtimeID", input, tenant, subAccountId)

		//then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not match any provider")
	})

	t.Run("should return error when no Gardener provider specified", func(t *testing.T) {
		// given
		uuidGeneratorMock := &mocks.UUIDGenerator{}
		uuidGeneratorMock.On("New").Return("id").Times(4)

		input := gqlschema.ProvisionRuntimeInput{
			ClusterConfig: &gqlschema.ClusterConfigInput{
				GardenerConfig: &gqlschema.GardenerConfigInput{},
			},
		}

		inputConverter := NewInputConverter(uuidGeneratorMock, nil, gardenerProject)

		//when
		_, err := inputConverter.ProvisioningInputToCluster("runtimeID", input, tenant, subAccountId)

		//then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "provider config not specified")
	})

}

func newInputConverterTester(uuidGenerator uuid.UUIDGenerator, releaseRepo release.Provider) *converter {
	return &converter{
		uuidGenerator: uuidGenerator,
		releaseRepo:   releaseRepo,
	}
}

func fixKymaGraphQLConfigInput() *gqlschema.KymaConfigInput {
	return &gqlschema.KymaConfigInput{
		Version: kymaVersion,
		Components: []*gqlschema.ComponentConfigurationInput{
			{
				Component: clusterEssentialsComponent,
				Namespace: kymaSystemNamespace,
			},
			{
				Component: coreComponent,
				Namespace: kymaSystemNamespace,
				Configuration: []*gqlschema.ConfigEntryInput{
					fixGQLConfigEntryInput("test.config.key", "value", util.BoolPtr(false)),
					fixGQLConfigEntryInput("test.config.key2", "value2", util.BoolPtr(false)),
				},
			},
			{
				Component: rafterComponent,
				Namespace: kymaSystemNamespace,
				SourceURL: util.StringPtr(rafterSourceURL),
			},
			{
				Component: applicationConnectorComponent,
				Namespace: kymaIntegrationNamespace,
				Configuration: []*gqlschema.ConfigEntryInput{
					fixGQLConfigEntryInput("test.config.key", "value", util.BoolPtr(false)),
					fixGQLConfigEntryInput("test.secret.key", "secretValue", util.BoolPtr(true)),
				},
			},
		},
		Configuration: []*gqlschema.ConfigEntryInput{
			fixGQLConfigEntryInput("global.config.key", "globalValue", util.BoolPtr(false)),
			fixGQLConfigEntryInput("global.config.key2", "globalValue2", util.BoolPtr(false)),
			fixGQLConfigEntryInput("global.secret.key", "globalSecretValue", util.BoolPtr(true)),
		},
	}
}

func fixGQLConfigEntryInput(key, val string, secret *bool) *gqlschema.ConfigEntryInput {
	return &gqlschema.ConfigEntryInput{
		Key:    key,
		Value:  val,
		Secret: secret,
	}
}
