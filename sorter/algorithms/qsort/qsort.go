package qsort

func quickSort(values []int ,left, right int){
	if left < right {
		temp := values[left]
		i,j := left ,right
		for i < j {
			for i < j  && values[j] >= temp {
				j--
			}
			if i < j {
				values[i] = values[j]
				i++
			}
			for i < j && values[i] < temp {
				i++
			}
			if i < j {
				values[j] = values[i]
				j--
			}
		}
		values[i] = temp
		quickSort(values,left,i-1)
		quickSort(values,i+1,right)
	}
}

func QuickSort(values []int){
	quickSort(values,0,len(values)-1)
}
