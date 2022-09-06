#include <iostream>
using namespace std;

template <class T>
class SharePtr {
 private:
  int *ref_;
  T *ptr_;

  void release() {
    if (ptr_ != nullptr && --*ref_ == 0) {
      cout << "relese" << *ptr_ << endl;
      delete ref_;
      delete ptr_;
    }
  }

 public:
  SharePtr() {}
  SharePtr(T *ptr) {
    ref_ = new int(1);
    ptr_ = ptr;
  }

  SharePtr(SharePtr<T> &other) {
    ref_ = other.ref_;
    ptr_ = other.ptr_;
    cout << "copy1" << endl;
  }

  T *operator->() { return ptr_; }

  T &operator*() { return *ptr_; }

  SharePtr<T> &operator=(SharePtr<T> &other) {
    if (&other == this) {
      return *this;
    }

    release();

    ++*other.ref_;
    ref_ = other.ref_;
    ptr_ = other.ptr_;
    cout << "copy2" << endl;
    return *this;
  }

  int getRef() {
    cout << "count:" << *ref_ << endl;
    return *ref_;
  }

  ~SharePtr() { release(); }
};

int test() {
  // SharePtr<string> p1(new string("aaa"));
  // p1.getRef();

  // SharePtr<string> p2(p1);
  // p2.getRef();

  // SharePtr<string> p3(new string("bbb"));
  // p3 = p1;

  // SharePtr<string> p4(nullptr);
  SharePtr<string> p1(new string("aaa"));
  SharePtr<string> p2(p1);
  p2 = p1;
  return 0;
}

int main() {
  test();
  return 0;
}
