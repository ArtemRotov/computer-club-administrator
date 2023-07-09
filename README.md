# computer-club-administrator
Требуется написать прототип системы, которая следит за работой компьютерного клуба, обрабатывает события и подсчитывает выручку за день и время занятости каждого стола.

Используемые технологии:
- Golang 1.20 (std only)
- Docker (для развертывания)

## Использование <a name="using"></a>
Входные данные представляют собой текстовый файл. Файл указывается первым аргументом при запуске программы.
Все тестовые файлы лежат в директории **/tests**, а ответы к ним в **/tests/answ**.
1. Для сборки docker-контейнера выполнить:
```shell
make build-docker
```
2. Для запуска docker-контейнера (с входным файлом tests/test.txt) выполнить:
```shell
make run-docker
```
3. Для запуска docker-контейнера с произвольным тестовым файлом из директории **/tests** выполнить:
```shell
docker run --rm -i computer-club-manager /app/tests/{test_name}.txt
```
4. Для запуска docker-контейнера с произвольным тестовым файлом из любой другой директории выполнить:
```shell
docker run --rm -i -v /{path_to_file}/{filename}.txt:/app/{filename}.txt computer-club-manager /app/{filename}.txt
```

## Unit-тестирование <a name="testing"></a>
Для запуска unit-тестов выполнить: 
```shell
make test
```