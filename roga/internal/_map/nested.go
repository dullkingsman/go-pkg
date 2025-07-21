package _map

type NestedMap[T comparable, H comparable, S any] map[T]map[H]S

func (m NestedMap[T, H, S]) GetValue(key1 T, key2 H) S {
	if m == nil || m[key1] == nil {
		var s S
		return s
	}

	return m[key1][key2]
}

func (m NestedMap[T, H, S]) GetMap(key T) map[H]S {
	if m == nil {
		return nil
	}

	return m[key]
}

func (m NestedMap[T, H, S]) SetValue(key1 T, key2 H, value S) {
	if m == nil {
		return
	}

	if m[key1] == nil {
		m[key1] = make(map[H]S)
	}

	m[key1][key2] = value
}

func (m NestedMap[T, H, S]) SetMap(key T, value map[H]S) {
	if m == nil {
		return
	}

	m[key] = value
}

func (m NestedMap[T, H, S]) RemoveValue(key1 T, key2 H) {
	if m == nil {
		return
	}

	delete(m[key1], key2)
}

func (m NestedMap[T, H, S]) RemoveMap(key T) {
	delete(m, key)
}
