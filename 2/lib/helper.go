package lib;

func splitByLength(s string, chunkSize int) []string {
    var chunks []string

    for i := 0; i < len(s); i += chunkSize {
        end := i + chunkSize
		end = min(end, len(s))
        chunks = append(chunks, s[i:end])
    }

    return chunks
}

func HasSequence(num string) bool {
	length := len(num)
	
	for i := 1; i <= length/2; i++ {
		if length % i != 0 { continue }

		units := splitByLength(num, i)
		allEqual := true
		for i := 0; i < len(units) - 1; i++ {
			if units[i] != units[i+1] {
				allEqual = false
				break;
			} 
		}
		if (allEqual) { return true }
	}
	
	return false
}

