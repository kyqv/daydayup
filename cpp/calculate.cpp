#include <iostream>
#include <stack>
#include <string>
#include <unordered_map>
#include <vector>

using namespace std;

bool isOperator(char ch) {
  return ch == '*' || ch == '/' || ch == '+' || ch == '-';
}

int calculate(char op, int num1, int num2) {
  if (op == '*')
    return num1 * num2;
  else if (op == '/')
    return num1 / num2;
  else if (op == '+')
    return num1 + num2;
  else if (op == '-')
    return num1 - num2;
  else
    return 0;
}

unordered_map<char, int> priority = {
    {'(', 4}, {'*', 3}, {'/', 3}, {'+', 1}, {'-', 1}};

bool startWithLeftBracket(char prev) {
  return prev == '{' || prev == '(' || prev == '[';
}

int main() {
  string expression;

  stack<int> s1;               //保存操作数
  stack<char> s2;              //保存操作符
  vector<int> postExpression;  //保存后缀表达式，操作符也按ASCII码存为int
  string num;  //从expression取出每个操作数，操作数>=10的时候必须这么做

  while (cin >> expression) {
    // 1. 中缀转后缀
    for (int i = 0; i < expression.length(); i++) {
      if (expression[i] == '(') {
        expression[i] = '(';
        s2.push(expression[i]);
      } else if (expression[i] == ')') {
        // 取出一个num
        if (num.length() > 0) {
          postExpression.push_back(stoi(num));
          num.clear();
        }
        //弹出(上面的全部操作符
        while (!s2.empty() && s2.top() != '(') {
          postExpression.push_back(s2.top());
          s2.pop();
        }
        if (!s2.empty()) {
          s2.pop();
        }
      } else {
        if (isOperator(expression[i])) {  //是操作符
          //取出操作数
          if (num.length() > 0) {
            postExpression.push_back(stoi(num));
            num.clear();
          }
          //处理类似这种 [-1 * 2 ], (+3 * 2),即 +/- 出现在（括号）表达式的开头
          if ((expression[i] == '+' || expression[i] == '-') &&
              ((i - 1 >= 0 && startWithLeftBracket(expression[i - 1])) ||
               (i == 0))) {
            postExpression.push_back(0);  //补0策略
          }
          //与栈顶操作符比较优先级，忽略左括号(
          while (!s2.empty() && s2.top() != '(' &&
                 priority[s2.top()] >= priority[expression[i]]) {
            postExpression.push_back(s2.top());
            s2.pop();
          }
          s2.push(expression[i]);
        } else {
          num.push_back(expression[i]);  //是数字，连续的数字算一个操作数
        }
      }
    }

    //如果操作数是最后一个，记得放入post expression
    if (num.length() > 0) {
      postExpression.push_back(stoi(num));
      num.clear();
    }
    //弹出s2剩余的全部操作符
    while (!s2.empty()) {
      postExpression.push_back(s2.top());
      s2.pop();
    }

    // 2. 后缀表达式求值
    for (int i = 0; i < postExpression.size(); i++) {
      if (!isOperator(postExpression[i])) {  //是操作数，存到栈里
        s1.push(postExpression[i]);
      } else {  //是操作符，取出2个操作数
        int l = 0, r = 0;
        if (!s1.empty()) {
          r = s1.top();
          s1.pop();
        }
        if (!s1.empty()) {
          l = s1.top();
          s1.pop();
        }

        int ret = calculate(postExpression[i], l, r);  //计算
        s1.push(ret);
      }
    }
    int res = s1.top();
    s1.pop();
    cout << res << endl;
  }
  return 0;
}
