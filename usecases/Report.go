package usecases


func (usecase *UsecaseImplemented) DailyReport()  {
    usecase.Repo.DailyReport2()
}

func (usecase *UsecaseImplemented) WeeklyReport()  {
    usecase.Repo.WeeklyReport2()
}

func (usecase *UsecaseImplemented) MonthlyReport()  {
    usecase.Repo.MonthlyReport2()
}

func (usecase *UsecaseImplemented) YearlyReport()  {
    usecase.Repo.YearlyReport2()
}