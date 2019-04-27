package verification

import (
	"rtiddsgo/verification/src"
)

const TopicName = "ComTwoA"

var ComTwoX = eb.Com_Two_X{
	Com_One_One_Z: eb.Com_One_One_Z{
		Com_One_One_X: eb.Com_One_One_X{
			W_A: true,
		},
		A: true,
		B: -1,
		C: 1,
		D: -2,
		E: 2,
		F: -3.0,
		G: -4.0,
		I: []bool{true},
		J: []int16{-5, -6},
		K: []string{"Seven", "Eight"},
		L: eb.Com_One_One_O{
			eb.Com_One_One_F{true, false},
			eb.Com_One_One_F{true, true},
			eb.Com_One_One_F{false, false},
		},
	},
	X_A: eb.Com_One_One_W{},
}
