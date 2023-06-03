package utils

// comparatorFunc is a comparator function that has to be understood in analogy to <(a, b), where a and b are numbers and "<" is a comparator function that returns true if a < b and false otherwise.
func MergeSort [K any] (value []K, comparatorFunc func(a K, b K) bool) []K {
    if len(value) < 2 {
        return value
    }

    floor := int(len(value) / 2)
    left_split, right_split := MergeSort(value[:floor], comparatorFunc), MergeSort(value[floor:], comparatorFunc)

    var merged = make([]K, 0, len(value))
    i, j := 0, 0
    for i != len(left_split) && j != len(right_split) {
        if comparatorFunc(left_split[i], right_split[j]) {
            merged = append(merged, left_split[i])
            i++
            continue
        }
        merged = append(merged, right_split[j])
        j++
    }
    for i != len(left_split) {
        merged = append(merged, left_split[i])
        i++
    }
    for j != len(right_split) {
        merged = append(merged, right_split[j])
        j++
    }
    return merged
}
