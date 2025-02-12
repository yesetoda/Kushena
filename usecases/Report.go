package usecases


func (usecase *UsecaseImplemented) DailyReport()  {
    usecase.Repo.DailyReport()
}

func (usecase *UsecaseImplemented) WeeklyReport()  {
    usecase.Repo.WeeklyReport()
}

func (usecase *UsecaseImplemented) MonthlyReport()  {
    usecase.Repo.MonthlyReport()
}

func (usecase *UsecaseImplemented) YearlyReport()  {
    usecase.Repo.YearlyReport()
}