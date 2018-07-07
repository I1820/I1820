package actions

func (as *ActionSuite) Test_AboutHandler() {
	res := as.JSON("/about").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "18.20 is leaving us")
}
