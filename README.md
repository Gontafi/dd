# File Archiver and Email Sender

Этот репозиторий представляет собой Go-проект, предоставляющий API для создания и чтения архивов файлов, а также отправки файлов на несколько электронных адресов.

## Функциональность

1. **Создание архива из файлов**  
   Эндпоинт: `/api/archive/files`  
   - Принимает несколько файлов в формате `form-data`.
   - Поддерживаемые форматы файлов: `.docx`, `.xml`, `.jpg`, `.png`.
   - Возвращает `.zip` архив.

2. **Получение информации об архиве**  
   Эндпоинт: `/api/archive/information`  
   - Возвращает информацию об архиве в формате JSON:
     ```json
     {
       "filename": "",
       "archive_size": 66912,
       "total_size": 3458826,
       "total_files": 4,
       "files": [
         {
           "file_path": "example/path/file1.txt",
           "size": 12345,
           "meme_type": "text/plain; charset=utf-8"
         },
         ...
       ]
     }
     ```
   - Поля:
     - `filename` — имя архива.
     - `archive_size` — размер архива в байтах.
     - `total_size` — общий размер всех файлов в байтах.
     - `total_files` — количество файлов в архиве.
     - `files` — информация о каждом файле.

3. **Отправка файла на email**  
   Эндпоинт: `/api/mail/file`  
   - Отправляет файл форматов `.pdf` или `.docx` на указанные email-адреса.
   - Email-адреса передаются в теле запроса.

## Использование

Для удобного тестирования запросов предоставлен **Postman Collection**. Импортируйте коллекцию в Postman и отправляйте запросы к API.

*Что в итоге? Интервьюеру было все понятно и он решил повторить три раза "я тебя понял". В итоге так и не вышел на связь. doodocs от будущего стажера просит многого, а дает буквально 0*
