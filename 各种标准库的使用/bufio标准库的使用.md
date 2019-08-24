Go语言的 bufio 包实现了带缓存的 I/O 操作, 使用起来还是很爽的, 主要涉及到下面一下函数:

```go
func NewReader(rd io.Reader) *Reader : 创建读缓冲区 
func NewWriter(w io.Writer) *Writer : 创建写缓冲区 

func (b *Reader) Peek(n int) ([]byte, error) : 返回缓冲区前n字节, 不移动读取指针 
func (b *Reader) Read(p []byte) (n int, err error) : 读取数据到p中 
func (b *Reader) ReadByte() (c byte, err error) : 读取一个字节数据 
func (b *Reader) UnreadRune() error : 将最后读取的一个字节数据设置为未读, 下次仍然可以读取 
func (b *Reader) Buffered() int : 缓冲区中缓冲的还没有读取的数据 
func (b *Reader) ReadRune() (r rune, size int, err error) : 读取一个字符, 如中文字符”啊”可以直接读取 
func (b *Reader) UnreadRune() error : 设置最后一次读的Rune未读, 若最后一次不是ReadRune, 返回error 
func (b *Reader) ReadSlice(delim byte) (line []byte, err error) : 读取数据直到遇到delim 
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error) : 读取一行数据, 根据\n或者\r\n 
func (b *Reader) ReadBytes(delim byte) (line []byte, err error) : 读取delim之前的所有字节数据 
func (b *Reader) ReadString(delim byte) (line string, err error) : 读取delim之前的所有string数据

func (b *Writer) Flush() error : 刷新数据, 将缓冲区数据写入io writer 
func (b *Writer) Available() int : 写缓冲区可用的空间, 默认最大空间是4096 
func (b *Writer) WriteString(s string) (int, error) : 写入一个string 
func (b *Writer) WriteByte(c byte) error : 写入一个Byte 
func (b *Writer) WriteRune(r rune) (size int, err error) : 写入一个字符, 例如’你’或者’c’ 
func (b *Writer) Write(p []byte) (nn int, err error) : 写入一个字节数组 
func (b *Reader) WriteTo(w io.Writer) (n int64, err error) : WriteTo 实现了 io.WriterTo.
func (b *Writer) ReadFrom(r io.Reader) (n int64, err error) : ReadFrom 实现了 io.ReaderFrom. 
```

# bufio.NewReader

有时候我们想完整获取输入的内容，而输入的内容可能包含空格，这种情况下可以使用`bufio`包来实现.

```go
func bufioDemo() {
	reader := bufio.NewReader(os.Stdin) // 从标准输入生成读对象
	fmt.Print("请输入内容：")
	text, _ := reader.ReadString('\n') // 读到换行
	text = strings.TrimSpace(text)
	fmt.Printf("%#v\n", text)
}
```

