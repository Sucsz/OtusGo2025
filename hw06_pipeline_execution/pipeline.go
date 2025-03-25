package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	ch := orDone(done, in)
	for _, stage := range stages {
		ch = orDone(done, stage(ch))
	}
	return ch
}

func orDone(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			// Если канал done закрыт, начинаем слив оставшихся данных из in,
			// чтобы разблокировать потенциально ожидающие горутины на отправку.
			//nolint: revive
			for range in {
			}
		}()
		defer close(out)

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}
