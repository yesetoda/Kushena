package usecases


func (usecase *UsecaseImplemented) DailyReport()  ([]byte, error) {
    return usecase.Repo.DailyReport()
}

func (usecase *UsecaseImplemented) WeeklyReport()  ([]byte, error) {
    return usecase.Repo.WeeklyReport()
}

func (usecase *UsecaseImplemented) MonthlyReport()  ([]byte, error) {
    return usecase.Repo.MonthlyReport()
}

func (usecase *UsecaseImplemented) YearlyReport()  ([]byte, error) {
    return usecase.Repo.YearlyReport()
}