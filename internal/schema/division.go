package schema

import (
	"github.com/jinzhu/copier"

	"github.com/opentreehole/backend/internal/model"
)

type DivisionCreateRequest struct {
	// 分区名称: 树洞、评教等等
	Name string `json:"name"`

	// 分区详情：前端暂时不用
	Description string `json:"description"`
}

type DivisionModifyRequest struct {
	DivisionCreateRequest

	// TODO: 置顶的树洞 id
	Pinned []int `json:"pinned"`
}

// DivisionDeleteRequest Admin only
type DivisionDeleteRequest struct {
	// ID of the target division that all the deleted division's holes will be moved to
	To int `json:"to" default:"1"`
}

type DivisionResponse struct {
	// 新版 id
	ID int `json:"id"`

	// 旧版 id
	DivisionID int `json:"division_id"`

	// 分区名称: 树洞、评教等等
	Name string `json:"name"`

	// 分区详情：前端暂时不用
	Description string `json:"description"`

	// TODO: 置顶的树洞
	Pinned []struct{} `json:"pinned"`
}

func (d *DivisionResponse) FromModel(m *model.Division, pinned []struct{}) *DivisionResponse {
	err := copier.Copy(d, m)
	if err != nil {
		panic(err)
	}

	d.Pinned = pinned
	return d
}
