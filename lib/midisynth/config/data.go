package config

/*
{
	"instruments": {
		"1": {
			"wave": {
				"sine": {}
			},
			"filters": [
				{
					"adsr": {
						"releaseLen": 1.0
					}
				},
				{
					"delay": {
						"interval": 0.5,
						"times": 4,
						"fadeOut": 0.5
					}
				}
			]
		}
	}
}

*/

type Data struct {
	Instruments map[string]Instrument `json:"instruments"`
}

type Instrument struct {
	Wave    string   `json:"wave"`
	Filters []Filter `json:"filters"`
}

type Filter map[string]map[string]interface{}
