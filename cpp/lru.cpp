#include <iostream>
#include <list>
#include <map>
#include <string>
#include <unordered_map>
#include <vector>
using namespace std;

class LRUCache {
  list<pair<int, int>> cache;
  unordered_map<int, list<pair<int, int>>::iterator> hash;
  int size;

 public:
  LRUCache(int capacity) : size(capacity) {}

  void print() {
    cout << "hash size:" << hash.size() << endl;
    for (list<pair<int, int>>::iterator it = cache.begin(); it != cache.end();
         it++) {
      cout << it->first << ":" << it->second << endl;
    }

    cout << "-------" << endl;
  }

  int get(int key) {
    auto it = hash.find(key);
    if (it == hash.end()) {
      return -1;
    }

    auto key_value = *it->second;
    //删除，并放到最前面去
    cache.erase(it->second);
    cache.push_front(key_value);
    // hash重置下
    hash[key] = cache.begin();
    return key_value.second;
  }

  void put(int key, int value) {
    auto it = hash.find(key);
    if (it != hash.end()) {
      auto key_value = *it->second;
      //删除，并放到最前面去
      cache.erase(it->second);
      cache.push_front(key_value);
      // hash重置下
      hash[key] = cache.begin();
      return;
    }

    //新的插入
    cache.push_front(make_pair(key, value));
    hash[key] = cache.begin();
    //如果超出了淘汰旧的
    if (cache.size() > size) {
      hash.erase(cache.back().first);
      cache.pop_back();
    }
  }
};

int main() {
  LRUCache cache(2);
  cache.put(1, 2);
  cache.put(2, 3);
  cache.print();
  cache.put(4, 34);
  cache.print();
  cache.put(5, 35);
  cache.print();
}
