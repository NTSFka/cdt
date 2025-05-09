
cmake_minimum_required(VERSION 3.22)

project(main)

add_executable(main main.cpp)

add_executable(test_test test.cpp)

enable_testing()
add_test(test_test test_test)
