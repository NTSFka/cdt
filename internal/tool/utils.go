package tool

// appendFiles appends files to the arguments if they are not nil or empty.
func appendFiles(args []string, filenames *[]string, all *string) []string {
	if filenames != nil && len(*filenames) > 0 {
		args = append(args, *filenames...)
	} else if all != nil {
		args = append(args, *all)
	}

	return args
}
