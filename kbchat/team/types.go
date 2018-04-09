package team

//Results from keybase team api call
type Results struct {
	Result struct {
		Members struct {
			Owners []struct {
				Uv struct {
					UID         string `json:"uid"`
					EldestSeqno int    `json:"eldestSeqno"`
				} `json:"uv"`
				Username string `json:"username"`
				FullName string `json:"fullName"`
				Active   bool   `json:"active"`
				NeedsPUK bool   `json:"needsPUK"`
			} `json:"owners"`
			Admins []struct {
				Uv struct {
					UID         string `json:"uid"`
					EldestSeqno int    `json:"eldestSeqno"`
				} `json:"uv"`
				Username string `json:"username"`
				FullName string `json:"fullName"`
				Active   bool   `json:"active"`
				NeedsPUK bool   `json:"needsPUK"`
			} `json:"admins"`
			Writers []struct {
				Uv struct {
					UID         string `json:"uid"`
					EldestSeqno int    `json:"eldestSeqno"`
				} `json:"uv"`
				Username string `json:"username"`
				FullName string `json:"fullName"`
				Active   bool   `json:"active"`
				NeedsPUK bool   `json:"needsPUK"`
			} `json:"writers"`
			Readers []interface{} `json:"readers"`
		} `json:"members"`
		KeyGeneration          int `json:"keyGeneration"`
		AnnotatedActiveInvites struct {
		} `json:"annotatedActiveInvites"`
		Settings struct {
			Open   bool `json:"open"`
			JoinAs int  `json:"joinAs"`
		} `json:"settings"`
		Showcase struct {
			IsShowcased       bool `json:"is_showcased"`
			AnyMemberShowcase bool `json:"any_member_showcase"`
		} `json:"showcase"`
	} `json:"result"`
}

//Output for Admins, Writers, and Owners
type Output []struct {
	Uv struct {
		UID         string `json:"uid"`
		EldestSeqno int    `json:"eldestSeqno"`
	} `json:"uv"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Active   bool   `json:"active"`
	NeedsPUK bool   `json:"needsPUK"`
}
