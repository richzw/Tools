package snippet

// https://mp.weixin.qq.com/s/rhe1hFgIfFDHZldGhnzfJQ
func getPercentageRoundsSubstring(percentage float64) string {
	symbols := "★★★★★★★★★★☆☆☆☆☆☆☆☆☆☆"
	offset := 10 - int(percentage*10.0)
	return symbols[offset*3 : (offset+10)*3]
}

// Bit set
/*
利用或操作`|`和空格将英文字符转换为小写
('a' | ' ') = 'a'
('A' | ' ') = 'a'
利用与操作`&`和下划线将英文字符转换为大写
('b' & '_') = 'B'
('B' & '_') = 'B'
利用异或操作`^`和空格进行英文字符大小写互换
('d' ^ ' ') = 'D'
('D' ^ ' ') = 'd'
判断两个数是否异号
int x = -1, y = 2;
boolean f = ((x ^ y) < 0); // true
int x = 3, y = 2;
boolean f = ((x ^ y) < 0); // false

利用求模（余数）的方式让数组看起来头尾相接形成一个环形，永远都走不完
arr[index % arr.length]

    // 在环形数组中转圈
    print(arr[index & (arr.length - 1)]);
    index++;
    // 在环形数组中转圈
    print(arr[index & (arr.length - 1)]);
    index--;

n & (n-1)这个操作在算法中比较常见，作用是消除数字n的二进制表示中的最后一个 1。
- 判断一个数是不是 2 的指数

一个数和它本身做异或运算结果为 0，即a ^ a = 0；一个数和 0 做异或运算的结果为它本身，即a ^ 0 = a。
- 136 题「只出现一次的数字」
- 268 题「丢失的数字 - 只要把所有的元素和索引做异或运算，成对儿的数字都会消为 0，只有这个落单的元素会剩下
*/
