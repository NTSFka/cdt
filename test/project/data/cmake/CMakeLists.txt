
cmake_minimum_required(VERSION 3.22)

project(main)

include(FetchContent)
FetchContent_Declare(
  fmt
  URL ${CMAKE_CURRENT_SOURCE_DIR}/fmt-11.2.0.zip
  DOWNLOAD_EXTRACT_TIMESTAMP Off
)
FetchContent_MakeAvailable(fmt)

add_executable(main main.cpp)
target_link_libraries(main PRIVATE fmt::fmt)

add_executable(main_test test.cpp)

enable_testing()
add_test(main_test main_test)
