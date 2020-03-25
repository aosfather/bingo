package service

//--系统级set管理实现
type UnsaftSet map[interface{}]interface{}

func (this UnsaftSet) Gets(start int, stop int) []interface{} {
	index := 0
	var result []interface{}
	pass := false
	if stop < 0 {
		pass = true
	}
	//轮询加入到结果集中
	adding := false
	for k, _ := range this {
		//起始点大于等于现在的位置，并且result长度没有到stop的点
		if start == index || adding {

			if stop < index && !pass {
				break
			}
			result = append(result, k)
			adding = true
		}
		index++
	}
	return result
}

func (this UnsaftSet) Add(v ...interface{}) {
	for _, i := range v {
		found := this[i]

		if found == nil {
			this[i] = new(interface{})
		}

	}

}

//交集
func (this UnsaftSet) Inter(s ...UnsaftSet) UnsaftSet {

	return nil
}

//取差集
func (this UnsaftSet) Diff(s ...*UnsaftSet) *UnsaftSet {
	return nil
}

//取并集
func (this *UnsaftSet) Union(s ...*UnsaftSet) *UnsaftSet {
	return nil
}

//集合管理
type SetStore struct {
}

func (this *SetStore) LRange(key string, start int, stop int) {

}

func (this *SetStore) Del(key string, fields ...string) {

}

func (this *SetStore) UnionStore(id string, keys ...string) {

}

func (this *SetStore) InterStore(id string, keys ...string) {

}

//交集
func (this *SetStore) Inter(keys ...string) UnsaftSet {

	return nil
}

//取不同
func (this *SetStore) Diff(keys ...string) UnsaftSet {
	return nil
}

//取并集
func (this *SetStore) Union(keys ...string) UnsaftSet {
	return nil
}
