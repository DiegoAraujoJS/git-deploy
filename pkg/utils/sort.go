package utils

func MergeSort [K any] (value []K, comparatorFunc func(a K, b K) bool) []K {
    if len(value) < 2 {
        return value
    }

    floor := int(len(value) / 2)
    left_split, right_split := MergeSort(value[:floor], comparatorFunc), MergeSort(value[floor:], comparatorFunc)

    var merged = make([]K, 0, len(value))
    i, j := 0, 0
    for i != len(left_split) && j != len(right_split) {
        if comparatorFunc(right_split[j], left_split[i]) {
            merged = append(merged, right_split[j])
            j++
            continue
        }
        merged = append(merged, left_split[i])
        i++
    }

    for ; i < len(left_split) ; i++ {
        merged = append(merged, left_split[i])
    }

    for ; j < len(right_split) ; j++ {
        merged = append(merged, right_split[j])
    }

    return merged
}
