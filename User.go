package main

type User struct {
	GUID string `json:"guid"`

	Rts []string `json:"rts"`
}


func (u *User) RemoveAt(i int) {
	u.Rts[i] = u.Rts[len(u.Rts)-1]
	u.Rts = u.Rts[:len(u.Rts)-1]
}
