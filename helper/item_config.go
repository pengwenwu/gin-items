package helper

func GetItemFields() (fieldsMap map[string][]string) {
	fieldsMap = make(map[string][]string)
	fieldsMap["base"] = []string{"item_id", "appkey", "channel", "name", "photo", "detail", "state", "last_dated", "dated"}
	fieldsMap["photos"] = []string{"photos.id", "photos.item_id", "photos.photo", "photos.sort", "photos.state", "photos.last_dated", "photos.dated"}
	fieldsMap["parameters"] = []string{"parameters.item_id", "parameters.parameters", "parameters.value", "parameters.sort", "parameters.state", "parameters.last_dated", "parameters.dated"}
	fieldsMap["skus"] = []string{"skus.sku_id", "skus.item_id", "skus.appkey", "skus.channel", "skus.item_name", "skus.sku_name", "skus.sku_photo", "skus.sku_code", "skus.bar_code", "skus.properties", "skus.state", "skus.last_dated", "skus.dated"}
	fieldsMap["props"] = []string{"props"}
	return
}
