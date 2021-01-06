package pp

import (
	"math"
)

func calculateManiaPP(od float64, stars float64, noteCount float64, score float64) float64 {

	//var strainStep := 400 * timeScale
	hit300Window := 34 + 3*(math.Min(10, math.Max(0, 10-od)))
	strainValue := math.Pow(5*math.Max(1, stars/0.2)-4, 2.2) / 135 * (1 + 0.1*math.Min(1, noteCount/1500))
	//if score <= 500000 {
	//	strainValue = 0
	//}
	if score <= 600000 {
		strainValue *= (score - 500000) / 100000 * 0.3
	} else if score <= 700000 {
		strainValue *= 0.3 + (score-600000)/100000*0.25
	} else if score <= 800000 {
		strainValue *= 0.55 + (score-700000)/100000*0.20
	} else if score <= 900000 {
		strainValue *= 0.75 + (score-800000)/100000*0.15
	} else {
		strainValue *= 0.9 + (score-900000)/100000*0.1
	}
	accValue := math.Max(0, 0.2-((hit300Window-34)*0.006667)) * strainValue * math.Pow(math.Max(0, score-960000)/40000, 1.1)
	ppMultiplier := 0.8

	return math.Pow(math.Pow(strainValue, 1.1)+math.Pow(accValue, 1.1), 1/1.1) * ppMultiplier
}
