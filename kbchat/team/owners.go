package team

//Owners of the team
func Owners(results Results) (Output, error) {
	var owners Output
	for _, owner := range results.Result.Members.Owners {
		owners = append(owners, owner)
	}
	return owners, nil
}
