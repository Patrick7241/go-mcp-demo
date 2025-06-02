package prompt

import "fmt"

func FirstPrompt(input string) string {
	return fmt.Sprintf(
		`请严格按以下规则执行任务：

1. 输入内容：
   - 用户输入：[用户自然语言问题]

2. 你需要：
   - 判断该输入属于以下哪一类：
     a. 闲聊或打招呼（如：你好，在吗，哈哈，天气真好）
     b. 寻求帮助（如：如何使用系统，这个功能怎么用）

3. 输出要求：
   - 仅返回分类标签（a 或 b）
   - 不要添加其他内容或解释

当前任务：
用户输入：%s

请按要求返回结果：`, input)
}
func SecondPrompt(input string) string {
	return fmt.Sprintf(
		`请严格按以下规则执行任务：

1. 输入内容：
   - 用户输入：[用户自然语言问题]

2. 你需要：
   - 判断该输入属于以下哪一类：
     a. 其他（如：非电商相关的咨询）
     b. 询问商品详情或电商平台运营相关信息（如用户、交易、物流、订单等）

3. 输出要求：
   - 仅返回分类标签（a 或 b）
   - 不要添加其他内容或解释

当前任务：
用户输入：%s

请按要求返回结果：`, input)
}

func ThirdPrompt(userInput string, sqlFileText string) string {
	return fmt.Sprintf(
		`请严格按以下规则执行任务：

1. 输入内容包含：
   - 用户问题：[问题描述]
   - SQL文件内容：[SQL语句列表]

2. 任务要求：
   - 判断用户问题是否能用 SQL 文件中现有的语句解答
   - 如果不能解答，**请仅返回字符 a**，不允许返回其他任何内容
   - 如果能解答，**请仅返回能直接回答问题的完整 SQL 语句**，必须完全匹配 SQL 文件中的语句，不能做任何修改、补充或解释

3. 严格禁止：
   - 不允许返回字符 b 或其他任意内容
   - 不允许返回部分语句或修改后的语句
   - 只能返回 a 或文件中原始 SQL 语句

示例：

不能解答时，返回：
a

能解答时，返回完整 SQL 语句，如：
SELECT * FROM users WHERE id = 100;

当前任务：
用户问题：%s
SQL文件内容：
%s

请严格按要求返回结果（只返回 a 或完全匹配的 SQL 语句）：`, userInput, sqlFileText)
}
