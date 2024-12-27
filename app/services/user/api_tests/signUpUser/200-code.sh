#!/bin/bash

# Параметры запроса
url="http://localhost:8080/user/sign-up"
method="POST"
headers=(
  "accept: application/json"
  "Content-Type: application/json"
)
data='{"username":"user007","email":"user@gmail.com","password":"Password123)","description":"string","interests":["Манхва","Футбол","Сашими","School"],"age":150}'

# Выполнение запроса
response=$(curl -X "$method" "$url" \
  -H "${headers[@]}" \
  -d "$data")

# Проверка ответа
if echo "$response" | jq -e '.userId' >/dev/null; then
  echo "Тест успешен!"
else
  echo "Тест неудачен."
fi