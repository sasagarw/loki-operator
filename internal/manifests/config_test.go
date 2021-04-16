package manifests_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/ViaQ/loki-operator/internal/manifests"
	"github.com/ViaQ/loki-operator/internal/manifests/internal"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

func TestConfigOptions_UserOptionsTakePrecedence(t *testing.T) {
	// regardless of what is provided by the default sizing parameters we should always prefer
	// the user-defined values. This creates an all-inclusive manifests.Options and then checks
	// that every value is present in the result
	opts := randomConfigOptions()

	res, err := manifests.ConfigOptions(opts)
	require.NoError(t, err)

	expected, err := json.Marshal(opts.Stack)
	require.NoError(t, err)

	actual, err := json.Marshal(res.Stack)
	require.NoError(t, err)

	assert.JSONEq(t, string(expected), string(actual))
}

func TestConfigOptions_AppliesStackSize(t *testing.T) {
	allSizes := []lokiv1beta1.LokiStackSizeType{
		lokiv1beta1.SizeOneXExtraSmall,
		lokiv1beta1.SizeOneXSmall,
		lokiv1beta1.SizeOneXMedium,
	}
	for _, size := range allSizes {
		res, err := manifests.ConfigOptions(manifests.Options{
			Name:      "aksjdfadsf",
			Namespace: "lkjsadfl",
			Image:     "docker.io/ubuntu",
			Stack:     lokiv1beta1.LokiStackSpec{
				Size: size,
			},
		})
		require.NoError(t, err)
		require.EqualValues(t, internal.StackSizeTable[size], res.Stack)
	}
}

func randomConfigOptions() manifests.Options {
	return manifests.Options{
		Name:      uuid.New().String(),
		Namespace: uuid.New().String(),
		Image:     uuid.New().String(),
		Stack: lokiv1beta1.LokiStackSpec{
			Size:              lokiv1beta1.SizeOneXExtraSmall,
			Storage:           lokiv1beta1.ObjectStorageSpec{},
			StorageClassName:  uuid.New().String(),
			ReplicationFactor: rand.Int31(),
			Limits: lokiv1beta1.LimitsSpec{
				Global: lokiv1beta1.LimitsTemplateSpec{
					IngestionLimits: lokiv1beta1.IngestionLimitSpec{
						IngestionRate:           rand.Int31(),
						IngestionBurstSize:      rand.Int31(),
						MaxLabelLength:          rand.Int31(),
						MaxLabelValueLength:     rand.Int31(),
						MaxLabelsPerSeries:      rand.Int31(),
						MaxStreamsPerUser:       rand.Int31(),
						MaxGlobalStreamsPerUser: rand.Int31(),
						MaxLineSize:             rand.Int31(),
					},
					QueryLimits: lokiv1beta1.QueryLimitSpec{
						MaxEntriesPerQuery: rand.Int31(),
						MaxChunksPerQuery:  rand.Int31(),
						MaxQuerySeries:     rand.Int31(),
					},
				},
				Tenants: map[string]lokiv1beta1.LimitsTemplateSpec{
					uuid.New().String(): {
						IngestionLimits: lokiv1beta1.IngestionLimitSpec{
							IngestionRate:           rand.Int31(),
							IngestionBurstSize:      rand.Int31(),
							MaxLabelLength:          rand.Int31(),
							MaxLabelValueLength:     rand.Int31(),
							MaxLabelsPerSeries:      rand.Int31(),
							MaxStreamsPerUser:       rand.Int31(),
							MaxGlobalStreamsPerUser: rand.Int31(),
							MaxLineSize:             rand.Int31(),
						},
						QueryLimits: lokiv1beta1.QueryLimitSpec{
							MaxEntriesPerQuery: rand.Int31(),
							MaxChunksPerQuery:  rand.Int31(),
							MaxQuerySeries:     rand.Int31(),
						},
					},
				},
			},
			Template: lokiv1beta1.LokiTemplateSpec{
				Compactor: lokiv1beta1.LokiComponentSpec{
					Replicas: rand.Int31(),
					NodeSelector: map[string]string{
						uuid.New().String(): uuid.New().String(),
					},
					Tolerations: []corev1.Toleration{
						{
							Key:               uuid.New().String(),
							Operator:          corev1.TolerationOpEqual,
							Value:             uuid.New().String(),
							Effect:            corev1.TaintEffectNoExecute,
							TolerationSeconds: pointer.Int64Ptr(rand.Int63()),
						},
					},
				},
				Distributor: lokiv1beta1.LokiComponentSpec{
					Replicas: rand.Int31(),
					NodeSelector: map[string]string{
						uuid.New().String(): uuid.New().String(),
					},
					Tolerations: []corev1.Toleration{
						{
							Key:               uuid.New().String(),
							Operator:          corev1.TolerationOpEqual,
							Value:             uuid.New().String(),
							Effect:            corev1.TaintEffectNoExecute,
							TolerationSeconds: pointer.Int64Ptr(rand.Int63()),
						},
					},
				},
				Ingester: lokiv1beta1.LokiComponentSpec{
					Replicas: rand.Int31(),
					NodeSelector: map[string]string{
						uuid.New().String(): uuid.New().String(),
					},
					Tolerations: []corev1.Toleration{
						{
							Key:               uuid.New().String(),
							Operator:          corev1.TolerationOpEqual,
							Value:             uuid.New().String(),
							Effect:            corev1.TaintEffectNoExecute,
							TolerationSeconds: pointer.Int64Ptr(rand.Int63()),
						},
					},
				},
				Querier: lokiv1beta1.LokiComponentSpec{
					Replicas: rand.Int31(),
					NodeSelector: map[string]string{
						uuid.New().String(): uuid.New().String(),
					},
					Tolerations: []corev1.Toleration{
						{
							Key:               uuid.New().String(),
							Operator:          corev1.TolerationOpEqual,
							Value:             uuid.New().String(),
							Effect:            corev1.TaintEffectNoExecute,
							TolerationSeconds: pointer.Int64Ptr(rand.Int63()),
						},
					},
				},
				QueryFrontend: lokiv1beta1.LokiComponentSpec{
					Replicas: rand.Int31(),
					NodeSelector: map[string]string{
						uuid.New().String(): uuid.New().String(),
					},
					Tolerations: []corev1.Toleration{
						{
							Key:               uuid.New().String(),
							Operator:          corev1.TolerationOpEqual,
							Value:             uuid.New().String(),
							Effect:            corev1.TaintEffectNoExecute,
							TolerationSeconds: pointer.Int64Ptr(rand.Int63()),
						},
					},
				},
			},
		},
	}
}
