package team

//Members of the team including all Admins, Owners, and Writers
func Members(results Results) (Output, error) {
	var output Output
	for _, admin := range results.Result.Members.Admins {
		output = append(output, admin)
	}
	for _, writer := range results.Result.Members.Writers {
		output = append(output, writer)
	}
	for _, owner := range results.Result.Members.Owners {
		output = append(output, owner)
	}
	return output, nil
}
