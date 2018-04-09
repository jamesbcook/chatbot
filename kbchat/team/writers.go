package team

//Writers of a team
func Writers(results Results) (Output, error) {
	var writers Output
	for _, writer := range results.Result.Members.Writers {
		writers = append(writers, writer)
	}
	return writers, nil
}
