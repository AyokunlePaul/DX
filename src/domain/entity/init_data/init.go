package init_data

import "DX/src/domain/entity/category"

type Data struct {
	Categories []category.Category `json:"categories"`
}
