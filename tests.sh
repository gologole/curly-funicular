#!/bin/bash

# Функция для загрузки файла в базу данных
upload_file() {
  local name=$1
  local data=$2
  local userid=$3
  local hashfile=$4

  grpcurl -plaintext -d "{
    \"name\": \"$name\",
    \"data\": \"$(echo -n "$data" | base64)\",
    \"userid\": $userid,
    \"hashfile\": \"$hashfile\"
  }" localhost:50051 api.Storage/UploadFile
}

# Функция для получения файла
get_file() {
  local name=$1
  local userid=$2

  grpcurl -plaintext -d "{
    \"name\": \"$name\",
    \"userid\": $userid
  }" localhost:50051 api.Storage/GetFile
}

# Функция для удаления файла
delete_file() {
  local name=$1
  local userid=$2

  grpcurl -plaintext -d "{
    \"name\": \"$name\",
    \"userid\": $userid
  }" localhost:50051 api.Storage/DeleteFile
}

# Функция для получения списка файлов
get_file_list() {
  local userid=$1

  grpcurl -plaintext -d "{
    \"userid\": $userid
  }" localhost:50051 api.Storage/GetFileList
}

# Загрузка файлов в базу данных
upload_file "file1.txt" "Content of file 1" 1 "hash1"
upload_file "file2.txt" "Content of file 2" 1 "hash2"
upload_file "file3.txt" "Content of file 3" 2 "hash3"

# Тестирование GetFile
echo "Testing GetFile:"
get_file "file1.txt" 1

# Тестирование DeleteFile
echo "Testing DeleteFile:"
delete_file "file2.txt" 1

# Тестирование GetFileList
echo "Testing GetFileList:"
get_file_list 1
