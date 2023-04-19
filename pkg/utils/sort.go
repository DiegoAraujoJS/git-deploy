package utils

func BubbleSort [K any] (value []K, comparatorFunc func(n K, m K) bool) []K {
    flag := false
    for i := 0; i < len(value) - 1 ; i++ {
        if comparatorFunc(value[i], value[i + 1]) {
            flag = true
            value[i], value[i + 1] = value[i + 1], value[i]
        }
    }
    if flag {
        return BubbleSort(value, comparatorFunc)
    }
    return value
}
