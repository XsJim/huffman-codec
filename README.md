# huffman-codec
golang 实现的 huffman 编解码程序

## 解释说明

### 使用说明

1. `-f string` `f 文件名,文件名（第二个文件名仅在检查时可用）`

2. `-m string` `m 方法`
    1. `encode` 对目标文件进行编码
    2. `decode` 对目标文件进行解码
    3. `check` 对目标文件进行 md5 检查，比对文件是否一致


### 编码-压缩流程
1. 读入文件到字节数组

2. 统计各个字节出现数量

3. 将出现数量不为 0 的字节提取并生成节点

4. 节点入队（优先队列）

5. 逐个取出，组合，再入队，直到队列中只剩下一个节点

6. 得到根节点，将树写入到输出文件

7. 将原文件byte长度写入到输出文件-防止后续解码文件时读出补位 bit

8. 计算出每个叶节点字节对应的编码，放入编码表

9. 再次读取文件，按照字符查找编码表的到编码，依次输入到文件中

10. 结束

### 解码-解压缩流程
1. 获取读取器
2. 通过 获取树（读取器）树根 来获得树根
3. 获取写入器
4. 通过 解码器（树根，读取器，写入器） 将文件从读取器读出并解码到写入器中

### 说明
1. 节点结构：
    1. 一定含有一个权值
    2. 一定含有左子树、右子树位置
    3. 对于叶节点来说，还需要一个字节来对应表示的字符
2. 树写入文件：
    1. 在这个树中，有两种不同类型的节点，叶子节点和非叶子节点，非叶子节点表示路径，叶子节点表示字符
    2. 表示的字节占用8位，对于携带的字符，约定以 8 位长来表示1一个字符，也就是一个字节
    3. 按照先序遍历存储树，如果当前节点是路径，则用 bit 0 表示，否则用 bit 1 表示，并在它后边紧跟着 8 bit 表示对应字符
3. 树读出文件；
    1. 首先读入一位标记，如果当前节点是树叶，则接着读 8 位，组成字符，并返回该节点
    2. 如果当前节点是非树叶，则先递归的构造它的左子树，然后递归的构造它的右子树，之后返回该节点

### 程序结构分析

1. 入口，入口程序负责读取控制台输入的变量，分析当前是要编码还是解码，将文件路径引入对应的方法
2. 编码
    1. 获取读取器
    2. 生成 byte 统计数组
    3. 获取优先队列-实现接口 `sort.Interface`, `heap.Interface` 使用 `heap` 进行相关操作
    4. 计数非 0 的统计元素组合成节点并入队
    5. 生成树
    6. 生成写入器
    7. 通过 树写入器（树，写入器） 将树写入文件
    8. 通过 获取编码表（树）（编码表） 获得编码表
    9. 读取器游标置 0
    10. 通过 编码器（编码表，写入器，读取器） 写入文件
3. 解码
    1. 获取读取器
    2. 通过 构造树（读取器）（树根） 获取树
    3. 获取写入器
    4. 通过 解码器（树，读取器，写入器） 输出解码文件
4. 读写
   1. 读取器
      1. 读入 1 bit
      2. 读入 1 byte
   2. 写入器
      1. 写入 1 bit
      2. 写入 1 byte
5. 树操作
   1. 新建一个树节点（节点权值、节点字符、左子、右子）
   2. 判定节点是否是树叶
