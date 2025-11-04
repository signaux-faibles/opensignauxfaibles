package parsing

import "testing"

func Test_ParsePInt(t *testing.T) {
	a, err1 := ParsePInt("9")
	b := 9

	if *a == b && err1 == nil {
		t.Log("ParsePInt: parse 9 -> OK")
	} else {
		t.Error("ParsePInt: parse 9 -> Fail: " + err1.Error())
	}

	c, err2 := ParsePInt("a")
	if *c == 0 && err2 != nil {
		t.Log("ParsePInt: n'est pas un entier -> OK")
	} else {
		t.Error("ParsePInt: n'est pas un entier -> Fail")
	}

	d, err2 := ParsePInt("")
	if d == nil && err2 == nil {
		t.Log("ParsePInt: une valeur vide -> OK")
	} else {
		t.Error("ParsePInt: une valeur vide -> Fail")
	}
}

func Test_ParsePFloat(t *testing.T) {
	a, err1 := ParsePFloat("0.62")
	b := float64(0.62)

	if *a == b && err1 == nil {
		t.Log("ParsePFloat: parse 0.62 -> OK")
	} else {
		t.Error("ParsePFloat: parse 0.62 -> Fail: " + err1.Error())
	}

	c, err2 := ParsePFloat("abcd")
	if *c == 0.0 && err2 != nil {
		t.Log("ParsePFloat: n'est pas un float -> OK")
	} else {
		t.Error("ParsePFloat: n'est pas un float -> Fail")
	}

	d, err2 := ParsePFloat("")
	if d == nil && err2 == nil {
		t.Log("ParsePFloat: une valeur vide -> OK")
	} else {
		t.Error("ParsePFloat: une valeur vide -> Fail")
	}
}
