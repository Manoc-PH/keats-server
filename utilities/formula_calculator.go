package utilities

import (
	"errors"
	"math"
	"server/constants"
	const_defaults "server/constants/defaults"
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

// Calculating Coins and XP reward and bonus when inputting an intake
func Calc_CnXP_On_Add_Intake(Cal_Added float32, Cal_Total float32, Cal_Max float32) (Coins int, XP int, Deductions int) {
	new_cal_total := Cal_Added + Cal_Total
	deductions := 0
	c := int((Cal_Added / Cal_Max) * float32(const_defaults.Default_Coin_Reward))
	x := int((Cal_Added / Cal_Max) * float32(const_defaults.Default_XP_Reward))
	// checking if they exceeded their plan
	// if they did, their coins are deducted while not touching their xp
	if new_cal_total > Cal_Max {
		// if it exceeds, only use the excess as the multiplier
		isAt100AfterAddedCal := (new_cal_total / Cal_Max) >= 1
		if isAt100AfterAddedCal {
			excess := new_cal_total - Cal_Max
			positive_added_cal := Cal_Added - excess
			deductions = int((excess / Cal_Max) * float32(const_defaults.Default_Coin_Reward))
			c := int((positive_added_cal / Cal_Max) * float32(const_defaults.Default_Coin_Reward))
			x := int((positive_added_cal / Cal_Max) * float32(const_defaults.Default_XP_Reward))
			return c, x, deductions
		}
		c = 0
		x = 0
		deductions = int((Cal_Added / Cal_Max) * float32(const_defaults.Default_Coin_Reward))
		return c, x, deductions
	}
	// Checks if they already have the bonus
	// if not, a bonus is added to the total coins and xp
	isAt90BeforeAddedCal := (Cal_Total / Cal_Max) >= 0.9
	isAt90AfterAddedCal := (new_cal_total / Cal_Max) >= 0.9
	if isAt90AfterAddedCal && !isAt90BeforeAddedCal {
		c = c + const_defaults.Default_Coin_Bonus
		x = x + const_defaults.Default_XP_Bonus
	}
	return c, x, deductions
}
func Calc_CnXP_On_Delete_Intake(Cal_Deleted float32, Cal_Total float32, Cal_Max float32) (Coins int, XP int) {
	new_cal_total := Cal_Deleted - Cal_Total
	c := int((Cal_Deleted/Cal_Max)*float32(const_defaults.Default_Coin_Reward)) * -1
	x := int((Cal_Deleted/Cal_Max)*float32(const_defaults.Default_XP_Reward)) * -1
	// Checks if they already have the bonus
	// if not, the bonus is removed from the total coins and xp
	isAt90BeforeDeletedCal := (Cal_Total/Cal_Max) >= 0.9 && (Cal_Total/Cal_Max) <= 1
	isAt90AfterDeletedCal := (new_cal_total/Cal_Max) >= 0.9 && (new_cal_total/Cal_Max) <= 1
	// we check if it was no longer 90 after deleting the cal
	// and if it was at 90 before deleting the cal
	if !isAt90AfterDeletedCal && isAt90BeforeDeletedCal {
		c = c - const_defaults.Default_Coin_Bonus
		x = x - const_defaults.Default_XP_Bonus
	}
	return c, x
}
