package utilities

import (
	"errors"
	"math"
	"server/constants"
	"time"
)

func Calculate_Calories(sex string, w int, h int, bmr_multiplier float32, cal_percentage int, bday time.Time) (int, error) {
	now := time.Now()
	age := math.Floor(now.Sub(bday).Hours() / 24 / 365)
	// For Males
	if sex == constants.Sex_Types.Male {
		// Mifflin-St Jeor Equation
		bmr := int(((10 * float64(w)) + (6.25 * float64(h)) - (5 * age) + 5))
		calories := int((bmr_multiplier * float32(bmr)) * (float32(cal_percentage) * 0.01))
		return calories, nil
	}
	// For Females
	if sex == constants.Sex_Types.Female {
		// Mifflin-St Jeor Equation
		bmr := int(((10 * float64(w)) + (6.25 * float64(h)) - (5 * age) - 161))
		calories := int((bmr_multiplier * float32(bmr)) * (float32(cal_percentage) * 0.01))
		return calories, nil
	}

	return 0, errors.New("an error occured in calculating calories")
}

func Calculate_Macros(Cal int, p int, c int, f int) (protein int, carbs int, fats int) {
	p_pcnt := int(float32(Cal) * (float32(p) * 0.01) / 4)
	c_pcnt := int(float32(Cal) * (float32(c) * 0.01) / 4)
	f_pcnt := int(float32(Cal) * (float32(f) * 0.01) / 9)
	return p_pcnt, c_pcnt, f_pcnt
}
