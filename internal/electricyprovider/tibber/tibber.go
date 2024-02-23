package tibber

import "github.com/tobiasjaster/HomeEnergyManagement/internal/electricyprovider"

type TibberElectricyProviderBuilder struct{}

func (b *TibberElectricyProviderBuilder) Build() electricyprovider.IElectricyProvider {
	return &TibberElectricyProvider{}
}

type TibberElectricyProvider struct{}
