package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	for _, stage := range stages {
		if stage == nil {
			continue
		}
		in = StreamData(stage, in, done)
	}
	return in
}

func StreamData(stage Stage, chanIn In, done In) Out {
	outChannel := make(Bi)
	go func() {
		isClose := false
		for {
			select {
			case <-done:
				if !isClose {
					isClose = true
					close(outChannel)
				}
			case v, ok := <-chanIn:
				if !ok {
					if !isClose {
						close(outChannel)
					}
					return
				}
				if !isClose {
					outChannel <- v
				}
			}
		}
	}()
	return stage(outChannel)
}
