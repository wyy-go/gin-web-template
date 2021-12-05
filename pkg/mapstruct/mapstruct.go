package mapstruct

import (
	"fmt"
	"github.com/tiaotiao/mapstruct"
	"reflect"
)

func Struct2Map(s interface{}) map[string]interface{} {
	return mapstruct.Struct2MapTag(s, "json")
}

func Map2Struct(m map[string]interface{}, st interface{}) error {
	return mapstruct.Map2StructTag(m, st, "json")
}

func Struct2MapSlice(v interface{}) []map[string]interface{} {
	iv := reflect.Indirect(reflect.ValueOf(v))
	if iv.IsNil() || !iv.IsValid() || iv.Type().Kind() != reflect.Slice {
		return make([]map[string]interface{}, 0)
	}

	n := iv.Len()
	result := make([]map[string]interface{}, n)
	for i := 0; i < n; i++ {
		result[i] = Struct2Map(iv.Index(i).Interface())
	}

	return result
}

func MapSlice2Tree(m []map[string]interface{}, idField, pidField string) []map[string]interface{} {
	result := []map[string]interface{}{}
	idMap := make(map[string]interface{})

	for _, v := range m {
		id := fmt.Sprint(v[idField])
		idMap[id] = v
	}

	for _, v := range m {
		pid := fmt.Sprint(v[pidField])
		if _, ok := idMap[pid]; !ok || pid == "" {
			result = append(result, v)
		} else {
			pv := idMap[pid].(map[string]interface{})
			if _, ok := pv["children"]; !ok {
				var n []map[string]interface{}
				n = append(n, v)
				pv["children"] = &n
			} else {
				nodes := pv["children"].(*[]map[string]interface{})
				*nodes = append(*nodes, v)
			}
		}
	}
	return result
}

//TreeStruct 树形结构体
type TreeStruct struct {
	Parent map[string]interface{} `json:"parent"`
	Child  []*TreeStruct          `json:"children"`
}

//MapSlice2TreeStruct ...
func MapSlice2TreeStruct(m []map[string]interface{}, idField, pidField string) (newtree []*TreeStruct) {
	var sTree []*TreeStruct
	for _, v := range m {
		t := TreeStruct{Parent: v}
		sTree = append(sTree, &t) //引用地址
	}
	mapTree := make(map[string]*TreeStruct) //构建map,id做键
	for _, v := range sTree {
		v.Child = make([]*TreeStruct, 0)
		id := fmt.Sprintf("%v", v.Parent[idField])
		mapTree[id] = v
	}
	for _, v := range sTree {
		pid := fmt.Sprintf("%v", v.Parent[pidField])
		if m, ok := mapTree[pid]; ok {
			m.Child = append(m.Child, v)
		} else {
			newtree = append(newtree, v)
		}
	}
	return newtree
}
