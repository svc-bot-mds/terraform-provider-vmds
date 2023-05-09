package shared

func rabbitMQRegionConverter(regionData []string) []struct {
	id       string
	name     string
	disabled bool
} {
	data := make(map[string]int)
	for _, region := range regionData {
		data[region] = len(region)
	}

	regionKeys := make([]string, 0, len(regionData))
	for _, region := range regionData {
		regionKeys = append(regionKeys, region)
	}

	regions := make([]struct {
		id       string
		name     string
		disabled bool
	}, 0, len(regionKeys))

	for _, region := range regionKeys {
		regions = append(regions, struct {
			id       string
			name     string
			disabled bool
		}{
			id:       region,
			name:     region,
			disabled: data[region] == 0,
		})
	}

	return regions
}
