package utils

func CloseEmptyStrcutChannel(ch chan struct{}) {
	select {
	case _, ok := <-ch:
		if !ok {
			return
		}
	default:
	}

	close(ch)
}
