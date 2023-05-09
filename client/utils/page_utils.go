package utils

import "github.com/svc-bot-mds/terraform-provider-vmds/client/model"

func GetNextPageInfo(page *model.PageInfo) *model.PageQuery {
	if page.TotalPages == 0 || page.Number == page.TotalPages-1 {
		return nil
	}
	pageQuery := model.PageQuery{
		Index: page.Number + 1,
		Size:  page.Size,
	}
	return &pageQuery
}
