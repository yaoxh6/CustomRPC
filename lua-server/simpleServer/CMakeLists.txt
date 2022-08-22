get_filename_component(CURRENT_SOURCE_DIR_NAME ${CMAKE_CURRENT_SOURCE_DIR} NAME)
set(LIB_NAME ${CURRENT_SOURCE_DIR_NAME}) # 工程名, 默认文件夹名

include_directories(${LUNA_INCLUDE_DIRS})

file(GLOB_RECURSE ALL_SOURCE_FILES
        *.h
        *.cpp
)

add_library(${LIB_NAME} STATIC ${ALL_SOURCE_FILES})
target_link_libraries(${LIB_NAME} luna)