package tekojar

func main() {
	t, err := New()
	if err != nil {
		panic(err)
	}
	//
	// go func() {
	// 	for log := range t.WatchService("simple-loop.jar") {
	// 		PrintLog("UI", 0, log)
	// 	}
	// }()

	t.StartAll()
}
