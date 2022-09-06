
#include <iostream>
#include <list>
#include <map>
#include <string>
#include <unordered_map>
#include <vector>
using namespace std;

// 1.暴力法
int strStr(string haystack, string needle) {
  int n = haystack.size(), m = needle.size();
  for (int i = 0; i + m <= n; i++) {
    bool flag = true;
    for (int j = 0; j < m; j++) {
      if (haystack[i + j] != needle[j]) {
        flag = false;
        break;
      }
    }
    if (flag) {
      return i;
    }
  }
  return -1;
}

// 2.kmp
int strStr2(string haystack, string needle) {
  int n = haystack.size(), m = needle.size();
  if (m == 0) {
    return 0;
  }
  // next点
  vector<int> pi(m);
  for (int i = 1, j = 0; i < m; i++) {
    while (j > 0 && needle[i] != needle[j]) {
      j = pi[j - 1];
    }
    if (needle[i] == needle[j]) {
      j++;
    }
    pi[i] = j;
  }
  //匹配
  for (int i = 0, j = 0; i < n; i++) {
    while (j > 0 && haystack[i] != needle[j]) {
      j = pi[j - 1];
    }
    if (haystack[i] == needle[j]) {
      j++;
    }
    if (j == m) {
      return i - m + 1;
    }
  }
  return -1;
}
