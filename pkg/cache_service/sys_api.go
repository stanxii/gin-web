package cache_service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"strings"
)

// 获取所有接口
func (s *RedisService) GetApis(req *request.ApiListRequestStruct) ([]models.SysApi, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetApis(req)
	}
	var err error
	list := make([]models.SysApi, 0)
	query := s.redis.Table(new(models.SysApi).TableName())
	method := strings.TrimSpace(req.Method)
	if method != "" {
		query = query.Where("method", "contains", method)
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		query = query.Where("path", "contains", path)
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		query = query.Where("category", "contains", category)
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator", "contains", creator)
	}
	// 查询条数
	err = query.Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.NoPagination {
			// 不使用分页
			err = query.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = query.Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 根据权限编号获取以api分类分组的权限接口
func (s *RedisService) GetAllApiGroupByCategoryByRoleId(roleId uint) ([]response.ApiGroupByCategoryResponseStruct, []uint, error) {
	if !global.Conf.System.UseRedis {
		// 不使用redis
		return s.mysql.GetAllApiGroupByCategoryByRoleId(roleId)
	}
	// 接口树
	tree := make([]response.ApiGroupByCategoryResponseStruct, 0)
	// 有权限访问的id列表
	accessIds := make([]uint, 0)
	allApi := make([]models.SysApi, 0)
	// 查询全部api
	err := s.redis.Find(&allApi).Error
	if err != nil {
		return tree, accessIds, err
	}
	// 查询当前角色拥有api访问权限的casbin规则
	casbins, err := s.GetCasbinListByRoleId(roleId)
	if err != nil {
		return tree, accessIds, err
	}

	// 通过分类进行分组归纳
	for _, api := range allApi {
		category := api.Category
		path := api.Path
		method := api.Method
		access := false
		for _, casbin := range casbins {
			// 该api有权限
			if path == casbin.V1 && method == casbin.V2 {
				access = true
				break
			}
		}
		// 加入权限集合
		if access {
			accessIds = append(accessIds, api.Id)
		}
		// 生成接口树
		existIndex := -1
		children := make([]response.ApiListResponseStruct, 0)
		for index, leaf := range tree {
			if leaf.Category == category {
				children = leaf.Children
				existIndex = index
				break
			}
		}
		// api结构转换
		var item response.ApiListResponseStruct
		utils.Struct2StructByJson(api, &item)
		item.Title = fmt.Sprintf("%s %s[%s]", item.Desc, item.Path, item.Method)
		children = append(children, item)
		if existIndex != -1 {
			// 更新元素
			tree[existIndex].Children = children
		} else {
			// 新增元素
			tree = append(tree, response.ApiGroupByCategoryResponseStruct{
				Title:    category + "分组",
				Category: category,
				Children: children,
			})
		}
	}
	return tree, accessIds, err
}
