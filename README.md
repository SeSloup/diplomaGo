## Спасибо за проверку! 
------
- Описание проекта:  
>Планировщик хранит задачи; каждая из них содержит дату дедлайна и заголовок с комментарием. Задачи могут повторяться по заданному правилу: например, ежегодно, через какое-то количество дней, в определённые дни месяца или недели. Если отметить такую задачу как выполненную, она переносится на следующую дату в соответствии с правилом. Обычные задачи при выполнении будут просто удаляться. 

- Список выполненных заданий со звёздочкой:  
    - Обработка дат (FullNextDate = true)
    - Поиск (Search = false )
    - Авторизация (реализовано)

- Инструкция по запуску: 

>**Для запуска из папки с репозиторием использовать команду**  
*docker -D build -t tododiplomas:sologub1 .*  
*docker stop sologubdiplomas*  
*docker rm sologubdiplomas*  
*docker run --name sologubdiplomas -p 7540:7540  tododiplomas:sologub1*


```sh
docker stop sologubdiplomas; docker rm sologubdiplomas; docker -D build -t tododiplomas:sologub1 . ;  docker run --name sologubdiplomas -p 7540:7540 tododiplomas:sologub1
```

Проверить ссылку: [localhost:7540](http://localhost:7540)  
Пароль для проверки: 123456

>**для запуска нескомпилированной версии локально**
```
go run main.go
```
---------------------
--------------------
#### **тесты не отработают т.к. будет треботаться логирование по паролю.*
для проверки надо в файде ./pkg/api/api.go заменить строчки
>	http.HandleFunc("/api/task", auth.Auth(taskHandler))  
	http.HandleFunc("/api/tasks", auth.Auth(tasksHandler))  
	http.HandleFunc("/api/task/done", auth.Auth(doneTaskHandler))  

на     
>	http.HandleFunc("/api/task", taskHandler)  
	http.HandleFunc("/api/tasks", tasksHandler)  
	http.HandleFunc("/api/task/done", doneTaskHandler)  

после проверки вернуть обратно
