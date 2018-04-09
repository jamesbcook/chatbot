package team

//Admins in team
func Admins(results Results) (Output, error) {
	var admins Output
	for _, admin := range results.Result.Members.Admins {
		admins = append(admins, admin)
	}
	return admins, nil
}
