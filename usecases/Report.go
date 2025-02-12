package usecases


func (usecase *UsecaseImplemented) Report()  {
    usecase.Repo.DailyReport()
}