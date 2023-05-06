package snippet

import "io"

func ConcurrencyWrtie(src io.Reader, dest [2]io.Writer) (err error) {
	errCh := make(chan error, 1)

	// 管道，主要是用来写、读流转化
	pr, pw := io.Pipe()
	// teeReader ，主要是用来 IO 流分叉
	wr := io.TeeReader(src, pw)

	// 并发写入
	go func() {
		var _err error
		defer func() {
			pr.CloseWithError(_err)
			errCh <- _err
		}()
		_, _err = io.Copy(dest[1], pr)
	}()

	defer func() {
		// TODO：异常处理
		pw.Close()
		_err := <-errCh
		_ = _err
	}()

	// 数据写入
	_, err = io.Copy(dest[0], wr)

	return err
}
