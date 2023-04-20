package utils

func MergeSort [K any] (value []K, comparatorFunc func(a K, b K) bool) []K {
    if len(value) < 2 {
        return value
    }

    floor := int(len(value) / 2)
    l1, l2 := MergeSort(value[:floor], comparatorFunc), MergeSort(value[floor:], comparatorFunc)

    var merged []K
    i, j := 0, 0
    for i != len(l1) && j != len(l2) {
        if comparatorFunc(l2[j], l1[i]) {
            merged = append(merged, l2[j])
            j++
        } else {
            merged = append(merged, l1[i])
            i++
        }
    }

    for ; i < len(l1) ; i++ {
        merged = append(merged, l1[i])
    }

    for ; j < len(l2) ; j++ {
        merged = append(merged, l2[j])
    }

    return merged
}
