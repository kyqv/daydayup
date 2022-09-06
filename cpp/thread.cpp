#include <functional>  // std::ref
#include <future>      // std::promise, std::future
#include <iostream>    // std::cout
#include <thread>      // std::thread
#include <vector>

using namespace std;

void WaitForMilkTea(future<int>& f_notice) {
  // 做点别的，比如逛街
  int notice = f_notice.get();  // 查看奶茶好了没
  cout << "收到通知，回来取奶茶" << notice << endl;
}

void MakeMilkTea(promise<int>& p_notice) {
  // 制作奶茶
  cout << "奶茶做好了，通知顾客" << endl;
  p_notice.set_value(1);
}

int main() {
  promise<int> p_notice;
  auto f_notice = p_notice.get_future();  // future与会通知顾客的promise相关联
  thread Customer(WaitForMilkTea, ref(f_notice));
  thread Waiter(MakeMilkTea, ref(p_notice));
  Waiter.join();
  Customer.join();

  vector<int> num({1, 2, 3});
  for (auto a : num) {
    cout << a << endl;
  }
}
