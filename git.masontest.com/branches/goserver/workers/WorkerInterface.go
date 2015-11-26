package workers

type Worker interface {
    DealInput(resStr string)
    DealProcess(resStr string)
}
