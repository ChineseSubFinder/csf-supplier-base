package settings

type ImdbInfoCenterConfig struct {
	BaseUrl         string   `json:"base_url"`
	Files           []string `json:"files"`
	ImportExistData bool     `json:"import_exist_data"`
	DBConnectConfig DBConnectConfig
}
