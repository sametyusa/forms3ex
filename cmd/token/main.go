// Candidates: please ignore the code in this file,
// it's out of scope for the exercise.
package main

import (
	"flag"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

func main() {
	var hotel = flag.Int("hotel", 123, "hotel ID")
	var admin = flag.Bool("admin", false, "Token has admin privileges: true/false")
	var key = flag.String("key", "SigningString", "Signing key. Defaults to match docker-compose")

	flag.Parse()

	claims := jwt.MapClaims{
		"sub":   "8790c514-73b6-400f-8f28-acc74d342a22",
		"name":  "H.A. Kerr",
		"hotel": hotel,
		"admin": admin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(*key))
	if err != nil {
		panic(err)
	}
	fmt.Print(tokenString)
}
