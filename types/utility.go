package types


func MatrixToSlice(numFloors int, numButtons int, getValue func(int, int) bool) [][]bool {

	slice := make([][]bool, numFloors)

	for f := 0; f < numFloors; f++ {
		slice[f] = make([]bool, numButtons)

		for b := 0; b < numButtons; b++ {
			slice[f][b] = getValue(f, b)
		}
	}

	return slice
}