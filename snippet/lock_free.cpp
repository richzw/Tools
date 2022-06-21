template<typename T, int N>
// 链表结点称之为chunk_t
struct chunk_t
{
    T values[N]; //每个chunk_t可以容纳N个T类型的元素，以后就以一个chunk_t为单位申请内存
    chunk_t *prev;
    chunk_t *next;
};

//  Adds an element to the back end of the queue.
inline void push()
{
   back_chunk = end_chunk;
   back_pos = end_pos; //
   if (++end_pos != N) //end_pos!=N表明这个chunk节点还没有满
        return;

    chunk_t *sc = spare_chunk.xchg(NULL); // 为什么设置为NULL？ 因为如果把之前值取出来了则没有spare chunk了，所以设置为NULL
    if (sc)                               // 如果有spare chunk则继续复用它
    {
        end_chunk->next = sc;
        sc->prev = end_chunk;
    }
    else // 没有则重新分配
    {
        // static int s_cout = 0;
        // printf("s_cout:%d\n", ++s_cout);
        end_chunk->next = (chunk_t *)malloc(sizeof(chunk_t)); // 分配一个chunk
        alloc_assert(end_chunk->next);
        end_chunk->next->prev = end_chunk;
    }
    end_chunk = end_chunk->next;
    end_pos = 0;
}

//  Removes an element from the front end of the queue.
inline void pop()
{
    if (++begin_pos == N) // 删除满一个chunk才回收chunk
    {
        chunk_t *o = begin_chunk;
        begin_chunk = begin_chunk->next;
        begin_chunk->prev = NULL;
        begin_pos = 0;

        //  'o' has been more recently used than spare_chunk,
        //  so for cache reasons we'll get rid of the spare and
        //  use 'o' as the spare.
        chunk_t *cs = spare_chunk.xchg(o); //由于局部性原理，总是保存最新的空闲块而释放先前的空闲快
        free(cs);
    }
}

/  Initialises the pipe.
inline ypipe_t()
//  The destructor doesn't have to be virtual. It is mad virtual
//  just to keep ICC and code checking tools from complaining.

//  Flush all the completed items into the pipe. Returns false if
//  the reader thread is sleeping. In that case, caller is obliged to
//  wake the reader up before using the pipe again.
// 刷新所有已经完成的数据到管道，返回false意味着读线程在休眠，在这种情况下调用者需要唤醒读线程。
// 批量刷新的机制， 写入批量后唤醒读线程；
// 反悔机制 unwrite
inline bool flush()
{
    //  If there are no un-flushed items, do nothing.
    if (w == f) // 不需要刷新，即是还没有新元素加入
        return true;

    //  Try to set 'c' to 'f'.
    // read时如果没有数据可以读取则c的值会被置为NULL
    if (c.cas(w, f) != w) // 尝试将c设置为f，即是准备更新w的位置
    {

        //  Compare-and-swap was unseccessful because 'c' is NULL.
        //  This means that the reader is asleep. Therefore we don't
        //  care about thread-safeness and update c in non-atomic
        //  manner. We'll return false to let the caller know
        //  that reader is sleeping.
        c.set(f); // 更新为新的f位置
        w = f;
        return false; //线程看到flush返回false之后会发送一个消息给读线程，这需要写业务去做处理
    }
    else  // 读端还有数据可读取
    {
        //  Reader is alive. Nothing special to do now. Just move
        //  the 'first un-flushed item' pointer to 'f'.
        w = f;             // 更新f的位置
        return true;
    }
}
//  Write an item to the pipe.  Don't flush it yet. If incomplete is
//  set to true the item is assumed to be continued by items
//  subsequently written to the pipe. Incomplete items are neverflushed down the stream.
// 写入数据，incomplete参数表示写入是否还没完成，在没完成的时候不会修改flush指针，即这部分数据不会让读线程看到。
inline void write(const T &value_, bool incomplete_)
{
    //  Place the value to the queue, add new terminator element.
    queue.back() = value_;
    queue.push();

    //  Move the "flush up to here" poiter.
    if (!incomplete_)
    {
        f = &queue.back(); // 记录要刷新的位置
        // printf("1 f:%p, w:%p\n", f, w);
    }
    else
    {
        //  printf("0 f:%p, w:%p\n", f, w);
    }
}

inline bool check_read()
{
    //  Was the value prefetched already? If so, return.
    if (&queue.front() != r && r) //判断是否在前几次调用read函数时已经预取数据了return true;
        return true;

    //  There's no prefetched value, so let us prefetch more values.
    //  Prefetching is to simply retrieve the
    //  pointer from c in atomic fashion. If there are no
    //  items to prefetch, set c to NULL (using compare-and-swap).
    // 两种情况
    // 1. 如果c值和queue.front()， 返回c值并将c值置为NULL，此时没有数据可读
    // 2. 如果c值和queue.front()， 返回c值，此时可能有数据度的去
    r = c.cas(&queue.front(), NULL); //尝试预取数据

    //  If there are no elements prefetched, exit.
    //  During pipe's lifetime r should never be NULL, however,
    //  it can happen during pipe shutdown when items are being deallocated.
    if (&queue.front() == r || !r) //判断是否成功预取数据
        return false;

    //  There was at least one value prefetched.
    return true;
}

//  Reads an item from the pipe. Returns false if there is no value.
//  available.
inline bool read(T *value_)
{
    //  Try to prefetch a value.
    if (!check_read())
        return false;

    //  There was at least one value prefetched.
    //  Return it to the caller.
    *value_ = queue.front();
    queue.pop();
    return true;
}