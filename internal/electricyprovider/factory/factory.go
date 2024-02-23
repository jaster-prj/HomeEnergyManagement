package factory

import (
	"github.com/tobiasjaster/HomeEnergyManagement/internal/electricyprovider"
	"github.com/tobiasjaster/HomeEnergyManagement/internal/electricyprovider/awattar"
	"github.com/tobiasjaster/HomeEnergyManagement/internal/electricyprovider/tibber"
)

var (
	electricyProviderMap = map[string]electricyprovider.IElectricyProviderBuilder{
		"awattar": &awattar.AwattarElectricyProviderBuilder{},
		"tibber":  &tibber.TibberElectricyProviderBuilder{},
	}
)

func GetElectricyProviderFactory(provider string) electricyprovider.IElectricyProviderBuilder {
	builder, ok := electricyProviderMap[provider]
	if !ok {
		return nil
	}
	return builder
}
